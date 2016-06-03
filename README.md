This is a volume plugin for Docker, allowing Docker containers
to enjoy persistent volumes residing on ploop, either locally
or on the distributed Virtuozzo Storage file system.

## Prerequisites

This plugin uses on the following software:
* [goploop](https://github.com/kolyshkin/goploop) (or [goploop-cli](https://github.com/kolyshkin/goploop-cli))
* [Docker volume plugin helper](https://github.com/docker/go-plugins-helpers/tree/master/volume)

For Virtuozzo Storage and/or ploop, you need to have [Virtuozzo](https://virtuozzo.com/) or [OpenVZ](https://openvz.org/) installed and running. As this is a plugin to Docker, naturally, you should have [Docker](https://docker.com/) up and running.

## Installation

This guide assumes you are using a recent version of Virtuozzo or OpenVZ, and have Docker up and running.

### From RPM repo

```
cd /etc/yum.repos.d/
wget https://goo.gl/9N6lfJ
yum install docker-volume-ploop
```

### From source

First, you need to have ```ploop-devel``` package installed:

```yum install ploop-devel```

Next, you need to have Go installed, and GOPATH environment variable set:

```
yum install golang git
echo 'export GOPATH=$HOME/go' >> ~/.bash_profile
echo 'PATH=$GOPATH/bin:$PATH' >> ~/.bash_profile
. ~/.bash_profile
```
 
 Finally, get the plugin:
 
```go get github.com/virtuozzo/docker-volume-ploop```

Generally, you don't have to install the dependencies, as ```go get``` will do it for you.

## Usage

You need to have this plugin started before starting docker daemon.
For available options, see

```docker-volume-ploop -help```

Most important, you need to provide a path where the plugin will store
its volumes. For example:

```docker-volume-ploop -home /some/path```

Next, you need to create a new volume. Example:

```docker volume create -d ploop -o size=512G -name MyFirstVolume```

Finally, run a container with the volume:

```docker run -it -v VOLUME:/MOUNT alpine /bin/ash```

Here ```VOLUME``` is the volume name, and ```MOUNT``` is the path under which
the volume will be available inside a container.

See ```man docker volume``` for other volume operations. For example, to list existing volumes:
 
 ```docker volume ls```
 
## Troubleshooting
 
### Docker with Virtuozzo/OpenVZ kernel

**For Docker to work, you need to make sure conntracks are enabled on the host.** In case it's not done, docker might complain like this:

```Error starting daemon: Error initializing network controller: error obtaining controller instance: failed to create NAT chain: iptables failed: iptables --wait -t nat -N DOCKER: iptables v1.4.21: can't initialize iptables table `nat': Table does not exist (do you need to insmod?)\nPerhaps iptables or your kernel needs to be upgraded.\n (exit status 3)```

To fix, edit ```/etc/modprobe.d/parallels.conf``` (or ```/etc/modprobe.d/openvz.conf```) to look like this:

```options nf_conntrack ip_conntrack_disable_ve0=0```

In other words, the value should be set to 0. After making the change, reboot the machine.

## Miscellaneous ploop operations

The following is the quick introduction of what operations can be performed with ploop images. For more detailed information about ploop, see [openvz.org/Ploop](https://openvz.org/Ploop).

Use ```ploop``` command line tool, and refer to an image by path to ```DiskDescriptor.xml``` file. This driver creates images under ```img``` subdirectory of its home. So, to use the following commands, you need to ```cd``` to the image directory, for example:

```cd /pcs/img/MyFirstVol/```

### Snapshots

To create a snapshot:

```ploop snapshot DiskDescriptor.xml```

To list snapshots:

```ploop snapshot-list DiskDescriptor.xml```
 
To delete a snapshot:

```ploop snapshot-delete -u UUID DiskDescriptor.xml```
 
To mount a snapshot (read-only):

```ploop mount -r -u UUID -m MOUNT_POINT DiskDescriptor.xml```
 
### Resizing
 
To resize an image (can be done while it's running):

```ploop resize -s SIZE DiskDescriptor.xml```
 
### Checking

In case something is wrong (ploop image can't be mounted etc.), you might want to check it.

```ploop check DiskDescriptor.xml```

If you want to run fsck on an inner filesystem, you can use the following command:

```ploop mount -F DiskDescriptor.xml```

Don't forget to unmount it:

```ploop umount DiskDescriptor.xml```

## Licensing

This software is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/kolyshkin/docker-volume-ploop/blob/master/LICENSE)
for the full license text.
