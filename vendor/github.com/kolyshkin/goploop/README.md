# goploop [![GoDoc](https://godoc.org/github.com/kolyshkin/goploop?status.png)](https://godoc.org/github.com/kolyshkin/goploop)

This is a Go wrapper for [libploop](https://github.com/kolyshkin/ploop/tree/master/lib),
a C library to manage ploop.

## What is ploop?

Ploop is a loopback block device (a.k.a. "filesystem in a file"), not unlike [loop](https://en.wikipedia.org/wiki/Loop_device) but with better performance
and more features, including:

* thin provisioning (image grows on demand)
* dynamic resize (both grow and shrink)
* instant online snapshots
* online snapshot merge
* optimized image migration with write tracker (ploop copy)

Ploop is implemented in the kernel and is currently available in OpenVZ RHEL6 and RHEL7 based kernels. For more information about ploop, see [openvz.org/Ploop](https://openvz.org/Ploop).

## Prerequisites

You need to have
* ext4 formatted partition (note RHEL/CentOS 7 installer uses xfs by default, that won't work!)
* ploop-enabled kernel installed
* ploop kernel modules loaded
* ploop-lib and ploop-devel packages installed

Currently, all the above comes with OpenVZ, please see [openvz.org/Quick_installation](https://openvz.org/Quick_installation).
After installing OpenVZ, you might need to run:

    yum install ploop-devel

## Building

If you are going to build a binary that uses this package statically,
you need to add `static_build` build tag to your `go build` command,
such as:

   go build -tags static_build

## Usage

This package is used by Docker ploop graphdriver, see https://github.com/kolyshkin/docker/tree/ploop/daemon/graphdriver/ploop

For primitive examples of how to use the package, see [ploop_test.go](ploop_test.go).
