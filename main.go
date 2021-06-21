package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/golang/glog"
)

var (
	input                     = flag.String("input", "", "input Dockerfile, required")
	output                    = flag.String("output", "-", "output Dockerfile. Default to stdout.")
	outputUsedVersionToStdOut = flag.Bool("output_used_version_to_stdout", false, "Output used version to stdout. Isn't compatible with --output=-")

	versionPlaceholder = flag.String("version_placeholder", "VERSION", "Additionally replace this with actual caddy version.")

	caddyVersionRE = regexp.MustCompile(`caddy:(\d\.\d\.\d)`)
)

func main() {
	flag.Parse()

	if *outputUsedVersionToStdOut && *output == "-" {
		glog.Fatal("--output_used_version_to_stdout can't be set when --output=-: stdout will be jammed.")
	}

	b, err := ioutil.ReadFile(*input)
	if err != nil {
		glog.Fatalf("Read input file failed, err: %v", err)
	}

	Output(replaceVersions(b, *versionPlaceholder))
}

func Output(b []byte, version string) {
	var w io.WriteCloser
	if *output == "-" {
		w = os.Stdout
	} else {
		var err error
		w, err = prepareFile()
		if err != nil {
			glog.Fatalf("Failed to prepare file to write, err: %v", err)
		}
	}

	if _, err := w.Write(b); err != nil {
		glog.Fatalf("Failed to write, err: %v", err)
	}
	fmt.Println(version)
}

func prepareFile() (io.WriteCloser, error) {
	if err := os.MkdirAll(filepath.Dir(*output), os.ModePerm); err != nil {
		return nil, err
	}

	if _, err := os.Stat(*output); os.IsNotExist(err) {
		return os.Create(*output)
	} else {
		glog.Error(err)
	}

	return os.Open(*output)
}

func replaceVersions(in []byte, placeholder string) ([]byte, string) {
	var versions []string

	for _, m := range caddyVersionRE.FindAllSubmatch(in, -1) {
		versions = append(versions, string(m[1]))
	}

	glog.Infof("All available versions: \n%v", strings.Join(versions, "\n"))

	if len(versions) == 0 {
		glog.Info("No caddy versions found, do nothing.")
		return in, ""
	}

	sort.Strings(versions)
	biggest := versions[len(versions)-1]

	return bytes.ReplaceAll(
		caddyVersionRE.ReplaceAll(in, []byte("caddy:"+biggest)),
		[]byte(placeholder), []byte(biggest)), biggest

}
