// go gen-con - Go Generic Containers
// Copyright 2016 Aurélien Rainone. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/aurelien-rainone/go-gencon/containers"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var (
	container = flag.String("cont", "", "generic container type name; must be set")
	containee = flag.String("type", "", "containee type name (user type or builtin); must be set")
	name      = flag.String("name", "", "override generated container name; default 'ContaineeContainer' or 'containeeContainer' if containee is exported")
	output    = flag.String("output", "", "output file name; default srcdir/containee_container.go")
	pkg       = flag.String("pkg", "", "package name of the generated file; default to 'main' for cli usage or the containee package if go-gencon is called by `go generate`")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of go-gencon:\n")
	fmt.Fprintf(os.Stderr, "    go-gencon [flags] -type containee -cont container\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Available generic containers:\n")
	fmt.Fprintf(os.Stderr, " - BoundedStack\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttp://github.com/aurelien-rainone/go-gencon\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {

	log.SetFlags(0)
	log.SetPrefix("go-gencon: ")
	flag.Usage = Usage
	flag.Parse()

	if os.Getenv("GOPACKAGE") != "" {
		*pkg = os.Getenv("GOPACKAGE")
	} else if *pkg == "" {
		*pkg = "main"
	}

	if *container == "" || *container == "" {
		flag.Usage()
		os.Exit(2)
	}

	var (
		g    Generator // generator instance
		cter string    // type name of generate container
	)
	cter = *containee + *container
	if len(*name) > 0 {
		cter = *name
	}

	g.nfo = GeneratorInfo{
		Containee: *containee,
		Container: cter,
		Package:   *pkg,
	}
	g.nfo.Exported, g.nfo.Builtin = typeInfoFromString(*containee)
	g.contType = *container

	// Print the header and package clause.
	g.Printf("// This file has been generated by \"go-gencon\"; DO NOT EDIT\n")
	g.Printf("// command: \"go-gencon %s\"\n", strings.Join(os.Args[1:], " "))
	g.Printf("// Go Generic Containers\n")
	g.Printf("// For more information see http://github.com/aurelien-rainone/go-gencon\n")
	g.Printf("\n")

	// Format the output.
	src := g.format()

	// generate output filename
	err := ioutil.WriteFile(
		generateOutputFileName(*containee, *container, *output),
		src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func generateOutputFileName(containee, container, outputFlag string) (fileName string) {
	baseName := strings.ToLower(fmt.Sprintf("%s_%s.go", containee, container))
	if outputFlag == "" {
		fileName = filepath.Join(".", baseName)
	} else {
		if info, err := os.Stat(outputFlag); err == nil && info.IsDir() {
			// -output is a directory
			fileName = filepath.Join(outputFlag, baseName)
		} else {
			// -output is just a filename
			fileName = filepath.Join(".", outputFlag)
		}
	}
	return
}

func typeInfoFromString(typename string) (isExported, isBuiltin bool) {
	isBuiltin = false
	switch typename {
	case "uint8", "uint16", "uint32", "uint64":
	case "int8", "int16", "int32", "int64":
	case "float32", "float64", "complex64", "complex128":
	case "byte", "rune":
	case "uint", "int", "uintptr":
		isBuiltin = true
	}
	isExported = isFirstUpper(typename)
	return
}

func isFirstUpper(s string) bool {
	r, n := utf8.DecodeRuneInString(s)
	if n == 0 {
		// should never get there
		panic("Can't get first letter of an empty string...")
	}
	if r == utf8.RuneError {
		// should never get there either
		panic("Invalid encoding")
	}
	return unicode.IsUpper(r)
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf      bytes.Buffer  // Accumulated output.
	nfo      GeneratorInfo // Info needed to generate the file
	contType string        // container to generate
}

type GeneratorInfo struct {
	Containee, Container, Package string
	Exported                      bool // containee is an exported type
	Builtin                       bool // containee is a builtin type
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	// list supported containers
	containers := map[string]string{
		"boundedstack": containers.BoundedStack,
	}

	// create container template
	var (
		tmpl  *template.Template
		str   string
		err   error
		found bool
	)
	if str, found = containers[strings.ToLower(g.contType)]; !found {
		log.Fatalf("No template for container '%s'", g.contType)
	}
	if tmpl, err = template.New(g.contType).Parse(str); err != nil {
		panic(err)
	}
	err = tmpl.Execute(&g.buf, g.nfo)
	if err != nil {
		panic(err)
	}

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
