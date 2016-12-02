This document only describes how to install the plugin from source.
For main documentation, see [README.md](README.md).

== Dependencies ==

First, you need to have ```ploop-devel``` package installed:

```yum install ploop-devel```

Next, you need to have Go installed, and GOPATH environment variable set:

```
yum install golang git
echo 'export GOPATH=$HOME/go' >> ~/.bash_profile
echo 'PATH=$GOPATH/bin:$PATH' >> ~/.bash_profile
. ~/.bash_profile
```

== Installation ==

Get the plugin:
 
```go get github.com/virtuozzo/docker-volume-ploop```

This should install all the Go packages that the plugin needs, and the plugin itself.

Now, install the configuration files:
 
 ```cd $GOPATH/src/github.com/*/docker-volume-ploop && make install```
 
 == Next steps ==
 
 Follow on to [README.md, section Starting](README.md#starting)
