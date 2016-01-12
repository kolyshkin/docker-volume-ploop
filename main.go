package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/docker/go-units"
	"github.com/kolyshkin/goploop"
)

// Options and their default values
var (
	home = flag.String("home", "/pcs", "Base directory where volumes are created")
	size = flag.String("size", "16GB", "Default image size")
	mode = flag.String("mode", "expanded", "Default ploop image mode")
	clog = flag.Uint("clog", 0, "Cluster block log size in 512-byte sectors")
	help = flag.Bool("help", false, "Print usage information")
)

func usage(ret int) {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", path.Base(os.Args[0]))
	flag.PrintDefaults()

	os.Exit(ret)
}

func main() {
	flag.Parse()

	if *help {
		usage(0)
	}

	// Fill in the default volume options
	var opts volumeOptions

	sizeBytes, err := units.RAMInBytes(*size)
	if err != nil {
		logrus.Fatalf("Can't parse size %s: %s\n", *size, err)
	}
	opts.size = uint64(sizeBytes >> 10) // convert to KB

	opts.mode, err = ploop.ParseImageMode(*mode)
	if err != nil {
		logrus.Fatalf("Can't parse mode %s: %s\n", *mode, err)
	}

	opts.clog = *clog

	// Let's run!
	d := newPloopDriver(*home, &opts)
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("root", "ploop"))
}
