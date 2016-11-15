package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/kolyshkin/goploop"
)

// Options and their default values
var (
	home  = flag.String("home", "/pcs", "Base directory where volumes are created")
	scope = flag.String("scope", "auto", "Volumes scope (local or global)")
	size  = flag.String("size", "16GB", "Default image size")
	mode  = flag.String("mode", "expanded", "Default ploop image mode")
	clog  = flag.String("clog", "0", "Cluster block log size in 512-byte sectors")
	tier  = flag.String("tier", "-1", "Virtuozzo Storage tier (0 is fastest")
	help  = flag.Bool("help", false, "Print usage information")
	debug = flag.Bool("debug", false, "Be verbose")
	quiet = flag.Bool("quiet", false, "Be quiet (errors only, to stderr)")
)

func usage(ret int) {
	fmt.Printf("Usage: %s [options]\n", path.Base(os.Args[0]))
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

	if err := opts.setSize(*size); err != nil {
		logrus.Fatalf(err.Error())
	}
	if err := opts.setMode(*mode); err != nil {
		logrus.Fatalf(err.Error())
	}
	if err := opts.setCLog(*clog); err != nil {
		logrus.Fatalf(err.Error())
	}
	if err := opts.setTier(*tier); err != nil {
		logrus.Fatalf(err.Error())
	}
	if err := opts.setScope(*scope); err != nil {
		logrus.Fatalf(err.Error())
	}

	// Set log level
	if *debug {
		if *quiet {
			logrus.Fatalf("Flags 'debug' and 'quiet' are mutually exclusive")
		}
		logrus.SetLevel(logrus.DebugLevel)
		ploop.SetVerboseLevel(ploop.Timestamps)
		logrus.Debugf("Debug logging enabled")
	}
	if *quiet {
		logrus.SetOutput(os.Stderr)
		logrus.SetLevel(logrus.ErrorLevel)
		ploop.SetVerboseLevel(ploop.NoStdout)
	}

	// Let's run!
	d := newPloopDriver(*home, &opts)
	h := volume.NewHandler(d)
	e := h.ServeUnix("root", "ploop")
	if e != nil {
		logrus.Fatalf("Failed to initialize: %s", e)
	}
}
