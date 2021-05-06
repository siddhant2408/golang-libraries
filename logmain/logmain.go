// Package logmain provies log related utilities for a main package.
//
// Importing this package automatically configures the default logger.
// The environment variable LOG_SHOW_FILE (bool) allows to display the current file in the logs.
package logmain

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/errors"
)

const (
	defaultFlags   = log.Ldate | log.Ltime | log.Lmicroseconds
	showFileEnvVar = "LOG_SHOW_FILE"
)

func init() {
	fs := defaultFlags
	showFile, err := getShowFile()
	if err != nil {
		err = errors.Wrap(err, "get show file")
		panic(err)
	}
	if showFile {
		fs |= log.Llongfile
	}
	log.SetFlags(fs)
}

func getShowFile() (bool, error) {
	s, ok := os.LookupEnv(showFileEnvVar)
	if !ok {
		return false, nil
	}
	ok, err := strconv.ParseBool(s)
	if err != nil {
		return false, errors.Wrap(err, "parse bool")
	}
	return ok, nil
}

// Start logs application start.
func Start(version string, env envutils.Env) {
	log.Println("Start")
	log.Printf("Version: %s", version)
	log.Printf("Environment: %s", env)
	log.Printf("Go version: %s", runtime.Version())
	log.Printf("Go max procs: %d", runtime.GOMAXPROCS(0))
}
