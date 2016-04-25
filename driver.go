package main

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/kolyshkin/goploop"
)

/* Driver options:
 * - home path (required)
 * - debug level
 * - defaults (all optional)
 *   - volume size
 *   - ploop format (expanded/preallocated/raw)
 *   - cluster block log size in 512-byte sectors,
 *     (values are from 6 to 15, default is 11: 2^11 * 512 = 1 MB)
 *
 * Volume options (for description see above):
 * - size (optional)
 * - format
 * - cluster block size
 */

type volumeOptions struct {
	size uint64          // ploop image size, in kilobytes
	mode ploop.ImageMode // ploop image format (expanded/prealloc/raw)
	clog uint            // cluster block log size in 512-byte sectors
}

type mount struct {
	count  int32
	device string
}

type ploopDriver struct {
	home    string
	opts    volumeOptions
	mountsM sync.RWMutex
	mounts  map[string]*mount
}

func newPloopDriver(home string, opts *volumeOptions) ploopDriver {
	// home must exist
	_, err := os.Stat(home)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Fatalf("Error %s", err)
		} else {
			logrus.Fatalf("Unexpected error from stat(%s): %s", home, err)
		}
	}

	d := ploopDriver{
		home:   home,
		opts:   *opts,
		mounts: make(map[string]*mount),
	}

	// Make sure to create base paths we'll use
	err = os.MkdirAll(d.img(""), 0700)
	if err != nil {
		logrus.Fatalf("Error %s", err)
	}
	err = os.MkdirAll(d.mnt(""), 0700)
	if err != nil {
		logrus.Fatalf("Error %s", err)
	}

	return d
}

func (d ploopDriver) Create(r volume.Request) volume.Response {
	// check if it already exists
	dd := d.dd(r.Name)
	_, err := os.Stat(dd)
	if err == nil {
		// volume already exists
		return volume.Response{}
	}
	if !os.IsNotExist(err) {
		logrus.Errorf("Unexpected error from stat(): %s", err)
		return volume.Response{Err: err.Error()}
	}

	logrus.Debugf("Creating volume %s", r.Name)
	// Create containing directory
	dir := d.dir(r.Name)
	err = os.Mkdir(dir, 0700)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	// Create an image
	file := d.img(r.Name)
	cp := ploop.CreateParam{Size: d.opts.size, Mode: d.opts.mode, File: file, CLog: d.opts.clog, Flags: ploop.NoLazy}

	if err := ploop.Create(&cp); err != nil {
		logrus.Errorf("Can't create ploop image: %s", err)
		return volume.Response{Err: err.Error()}
	}

	// all went well
	return volume.Response{}
}

func (d ploopDriver) Remove(r volume.Request) volume.Response {
	logrus.Debugf("Removing volume %s", r.Name)

	/* The ploop image to be removed might be mounted.
	 * The question is, what is the more correct thing to do:
	 * 1. Auto-unmount and proceed
	 * 2. Reject removing mounted image
	 */
	p, err := ploop.Open(d.dd(r.Name))
	if err == nil {
		if m, _ := p.IsMounted(); m {
			//err := fmt.Error("Rejecting to remove mounted image %s\n", r.Name)
			logrus.Error(err)
			return volume.Response{Err: err.Error()}
			/*
				err = p.Umount()
				if err != nil && !ploop.IsNotMounted(err) {
					logrus.Errorf("Can't umount %s: %s", r.Name, err)
					return volume.Response{Err: err.Error()}
				}
			*/
		}
		p.Close()
	}

	// Proceed with removal
	err = os.RemoveAll(d.dir(r.Name))
	if err != nil {
		logrus.Error(err)
		return volume.Response{Err: err.Error()}
	}

	// all went well
	return volume.Response{}
}

func (d ploopDriver) Mount(r volume.Request) volume.Response {
	logrus.Debugf("Mounting volume %s", r.Name)

	p, err := ploop.Open(d.dd(r.Name))
	if err != nil {
		logrus.Errorf("Can't open ploop: %s", err)
		return volume.Response{Err: err.Error()}
	}
	defer p.Close()

	mnt := d.mnt(r.Name)
	err = os.Mkdir(mnt, 0700)
	if err != nil && !os.IsExist(err) {
		logrus.Errorf("Error %s", err)
		return volume.Response{Err: err.Error()}
	}

	mp := ploop.MountParam{Target: mnt}

	dev, err := p.Mount(&mp)
	if err != nil {
		logrus.Errorf("Can't mount ploop: %s", err)
		return volume.Response{Err: err.Error()}
	}
	logrus.Debugf("Mounted %s to %s (dev=%s)", r.Name, d.mnt(r.Name), dev)

	// all went well
	return volume.Response{Mountpoint: mnt}
}

func (d ploopDriver) Unmount(r volume.Request) volume.Response {
	logrus.Debugf("Unmounting volume %s", r.Name)

	p, err := ploop.Open(d.dd(r.Name))
	if err != nil {
		logrus.Errorf("Can't open ploop: %s", err)
		return volume.Response{Err: err.Error()}
	}
	defer p.Close()

	if m, _ := p.IsMounted(); !m {
		// not mounted, nothing to do
		return volume.Response{}
	}

	err = p.Umount()
	// ignore "is not mounted" error
	if err != nil && !ploop.IsNotMounted(err) {
		logrus.Errorf("Can't unmount ploop: %s", err)
		return volume.Response{Err: err.Error()}
	}

	// all went well
	return volume.Response{}
}

func (d ploopDriver) Get(r volume.Request) volume.Response {
	logrus.Debugf("Called Get(%s)", r.Name)

	exist, err := d.volExist(r.Name)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}
	if !exist {
		// no such volume
		return volume.Response{Err: "Can't find volume"}
	}

	// TODO: check if it's mounted
	return volume.Response{Volume: &volume.Volume{Name: r.Name, Mountpoint: d.mnt(r.Name)}}
}

func (d ploopDriver) List(r volume.Request) volume.Response {
	logrus.Debugf("Called List()")
	dir := d.dir("")

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Errorf("Can't list directory %s: %s", dir, err)
		return volume.Response{Err: err.Error()}
	}

	vols := make([]*volume.Volume, 0, len(files))

	for _, f := range files {
		if f.IsDir() {
			name := f.Name()
			// Check if DiskDescriptor.xml is there
			exist, _ := d.volExist(name)
			if !exist {
				continue
			}
			vol := &volume.Volume{
				Name:       name,
				Mountpoint: d.mnt(name),
			}
			vols = append(vols, vol)
		}
	}

	return volume.Response{Volumes: vols}
}

func (d ploopDriver) Path(r volume.Request) volume.Response {
	logrus.Debugf("Called Path (%s)", r.Name)

	exist, err := d.volExist(r.Name)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if !exist {
		return volume.Response{Err: "Can't find volume"}
	}

	// TODO: check if mounted?
	return volume.Response{Mountpoint: d.mnt(r.Name)}
}

// Check if a given volume exist
func (d ploopDriver) volExist(name string) (bool, error) {
	dd := d.dd(name)
	_, err := os.Stat(dd)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		// no such volume
		return false, nil
	} else {
		logrus.Errorf("Unexpected error from stat(%s): %s", dd, err)
		return false, err
	}
}
