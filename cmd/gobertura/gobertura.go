package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/nim4/gocover-cobertura/cobertura"
	"golang.org/x/tools/cover"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		flagInput  string
		flagOutput string
		flagSrc    string
		flagPkg    string
	)
	flag.StringVar(&flagInput, "in", "coverprofile.txt", "path of coverage profile")
	flag.StringVar(&flagOutput, "out", "coverage.xml", "output path")
	flag.StringVar(&flagSrc, "src", "", "go source folder(will use current working directory if not set)")
	flag.StringVar(&flagPkg, "pkg", "", "package import path(will use `go.mod` if not set)")
	flag.Parse()

	convert(flagSrc, flagPkg, flagInput, flagOutput)
}

func convert(src string, pgk string, in string, out string) {
	if pgk == "" {
		data, err := ioutil.ReadFile("go.mod")
		if err != nil {
			panic(err)
		}

		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "module ") {
				pgk = strings.TrimSpace(strings.TrimPrefix(line, "module ")) + "/"
			}
		}
	}

	if src == "" {
		var err error
		src, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	profiles, err := cover.ParseProfiles(in)
	if err != nil {
		panic(err)
	}

	coverage := cobertura.Coverage{
		PackagePath: pgk,
		Sources: []*cobertura.Source{
			{
				Path: src,
			},
		},
		Packages:  nil,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
	err = coverage.ParseProfiles(profiles)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	write(f, xml.Header)
	write(f, "<!DOCTYPE coverage SYSTEM \"http://cobertura.sourceforge.net/xml/coverage-04.dtd\">\n")

	encoder := xml.NewEncoder(f)
	encoder.Indent("", "\t")
	err = encoder.Encode(coverage)
	if err != nil {
		panic(err)
	}

	write(f, "\n")
}

func write(f *os.File, str string) {
	_, err := fmt.Fprintf(f, str)
	if err != nil {
		panic(err)
	}
}
