package main

import "path"

const (
	ddxml       = "DiskDescriptor.xml"
	imagePrefix = "root.hdd"
)

// Returns path to ploop image directory for given id
func (d *ploopDriver) dir(id string) string {
	// Assuming that id doesn't contain "/" characters
	return path.Join(d.home, "img", id)
}

// Returns path to ploop's DiskDescriptor.xml for given id
func (d *ploopDriver) dd(id string) string {
	return path.Join(d.dir(id), ddxml)
}

// Returns path to ploop's image for given id
func (d *ploopDriver) img(id string) string {
	return path.Join(d.dir(id), imagePrefix)
}

// Returns a mount point for given id
func (d *ploopDriver) mnt(id string) string {
	return path.Join(d.home, "mnt", id)
}
