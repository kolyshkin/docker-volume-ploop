This is a volume plugin for Docker, allowing Docker containers
to enjoy persistent volumes residing on ploop, either locally
or on the distributed Virtuozzo Storage file system.

## Prerequisites

This plugin relies of the following software:
* [goploop](https://github.com/kolyshkin/goploop) (or [goploop-cli](https://github.com/kolyshkin/goploop-cli))
* [Docker volume plugin helper](https://github.com/docker/go-plugins-helpers/tree/master/volume)

## Installation

The following assumes you are using a recent version of Virtuozzo or OpenVZ.

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
 
```go get github.com/kolyshkin/docker-volume-ploop```

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
 
## Licensing

This software is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/kolyshkin/docker-volume-ploop/blob/master/LICENSE)
for the full license text.
