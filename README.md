This is a volume plugin for Docker, allowing Docker containers
to enjoy persistent volumes residing on ploop, either local or
with distributed pstorage file system.

## Prerequisites

This plugin relies of the following software:
* [goploop](https://github.com/kolyshkin/goploop) (or [goploop-cli](https://github.com/kolyshkin/goploop-cli))
* [Docker volume plugin helper](https://github.com/docker/go-plugins-helpers/tree/master/volume)

## Using

You need to have this plugin started before starting docker daemon.
For available options, see ```./docker-volume-ploop -help```.

An example of running container with a ploop volume:

```docker run -it --volume-driver ploop -v VOLUME:/MOUNT alpine /bin/ash```

Here ```VOLUME``` is the volume name, and ```MOUNT``` is the path under which
the volume will be available inside a container.

## Licensing

This software is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/kolyshkin/docker-volume-ploop/blob/master/LICENSE)
for the full license text.
