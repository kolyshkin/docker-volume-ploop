package main

import (
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
	d := ploopDriver{
		home:   home,
		opts:   *opts,
		mounts: make(map[string]*mount),
	}

	return d
}

func (d ploopDriver) Create(r volume.Request) volume.Response {
	logrus.Debugf("Creating volume %s\n", r.Name)

	// check if it already exists
	dd := d.dd(r.Name)
	_, err := os.Stat(dd)
	if err == nil {
		// volume already exists
		return volume.Response{}
	}
	if !os.IsNotExist(err) {
		logrus.Errorf("Unexpected error from stat(): %s\n", err)
		return volume.Response{Err: err.Error()}
	}

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
	logrus.Debugf("Removing volume %s\n", r.Name)

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
	logrus.Debugf("Mounting volume %s\n", r.Name)

	p, err := ploop.Open(d.dd(r.Name))
	if err != nil {
		logrus.Errorf("Can't open ploop: %s\n", err)
		return volume.Response{Err: err.Error()}
	}
	defer p.Close()

	var mp ploop.MountParam
	mp.Target = d.mnt(r.Name)

	dev, err := p.Mount(&mp)
	if err != nil {
		logrus.Errorf("Can't mount ploop: %s\n", err)
		return volume.Response{Err: err.Error()}
	}
	logrus.Debugf("Mounted %s to %s (dev=%s)\n", r.Name, d.mnt(r.Name), dev)

	// all went well
	return volume.Response{}
}

func (d ploopDriver) Unmount(r volume.Request) volume.Response {
	logrus.Debugf("Unmounting volume %s\n", r.Name)

	p, err := ploop.Open(d.dd(r.Name))
	if err != nil {
		logrus.Errorf("Can't open ploop: %s\n", err)
		return volume.Response{Err: err.Error()}
	}
	defer p.Close()

	err = p.Umount()
	// ignore "is not mounted" error
	if err != nil && !ploop.IsNotMounted(err) {
		logrus.Errorf("Can't unmount ploop: %s\n", err)
		return volume.Response{Err: err.Error()}
	}

	// all went well
	return volume.Response{}
}

func (d ploopDriver) Path(r volume.Request) volume.Response {
	return volume.Response{Mountpoint: d.mnt(r.Name)}
}
