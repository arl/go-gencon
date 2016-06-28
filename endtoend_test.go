// go gen-con - Go Generic Containers
// Copyright 2016 Aurélien Rainone. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// go command is not available on android

// +build !android

package main

import (
	"fmt"
	//"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// This file contains a test that compiles and runs every container in
// ./containers dir against each test program in ./testdata, after generating a
// specific container. Containee type name is read from a special comment in
// the test file. The rule is that for each containers/cont.go and each
// testdata/file.go we run a command that is `go-gencon -cont Cont -type typename`
// and then compile and run the program. The resulting binary panics if
// CONTINUE HERE!!

// 1) for every file in testdata/file.go
// 2) for every file in containers/cont.go
// 3) parse file.go and find the comment starting with endtoend
// 4) add -cont Cont to the previously found command
// 5) generate the container described in main.go
// 6) compile and run main.go

func TestEndToEnd(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-gencon")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	// Create go-gencon in temporary directory.
	gogencon := filepath.Join(dir, "go-gencon.exe")
	err = run("go", "build", "-o", gogencon, "go-gencon.go")
	if err != nil {
		t.Fatalf("building go-gencon: %s", err)
	}

	// Read the containers directory.
	fd, err := os.Open("containers")
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	names, err := fd.Readdirnames(-1)
	if err != nil {
		t.Fatalf("Readdirnames: %s", err)
	}

	// Generate, compile, and run the test programs.
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			t.Errorf("%s is not a Go file", name)
			continue
		}
		// Get corresponding test file name
		testfile := fmt.Sprintf("testdata/%s_test.go", name[:len(name)-len(".go")])
		args := readGoGenConParams(t, testfile)

		// Names are known to be ASCII and long enough.
		typeName := fmt.Sprintf("%c%s", name[0]+'A'-'a', name[1:len(name)-len(".go")])
		fmt.Println("stringerCompileAndRun", dir, gogencon, typeName, name, args)
		goGenConCompileAndRun(t, dir, gogencon, typeName, name, args)
	}
}

func readGoGenConParams(t *testing.T, filename string) string {
	// parse and retrieve all the test file comments
	fs := token.NewFileSet()
	astFile, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parsing test file: %s: %s", filename, err)
	}

	// compile the regex that we'll use to find the special comment
	rgx := `endtoend_test: go-gencon (.*)`
	r, _ := regexp.Compile(rgx)
	for i := range astFile.Comments {
		// walk through the comments
		cl := (*astFile.Comments[i]).List
		for j := range cl {
			cmt := cl[j].Text
			if matches := r.FindAllStringSubmatch(cmt, -1); matches != nil {
				return string(matches[0][0])
			}
		}
	}
	// shouldn't be here
	t.Fatalf("No comment in %s contains the go-gencon special comment", filename)
	return ""
}

// TODO: continuer ici : LIRE les parametres de go-gencon dans le fichier de test
// pour ca il faut parser le fichier
//modifier gogenConcompileandrun  pour accepter la liste de parameteres a passer a gogencon

// goGenConCompileAndRun runs go-gencon for the named file and compiles and
// runs the target binary in directory dir. That binary will panic if the String method is incorrect.
func goGenConCompileAndRun(t *testing.T, dir, stringer, typeName, fileName, args string) {
	t.Logf("run: %s %s\n", fileName, typeName)
	source := filepath.Join(dir, fileName)
	err := copy(source, filepath.Join("testdata", fileName))
	if err != nil {
		t.Fatalf("copying file to temporary directory: %s", err)
	}
	stringSource := filepath.Join(dir, typeName+"_string.go")
	// Run stringer in temporary directory.
	err = run(stringer, "-type", typeName, "-output", stringSource, source)
	if err != nil {
		t.Fatal(err)
	}
	// Run the binary in the temporary directory.
	err = run("go", "run", stringSource, source)
	if err != nil {
		t.Fatal(err)
	}
}

// copy copies the from file to the to file.
func copy(to, from string) error {
	toFd, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFd.Close()
	fromFd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFd.Close()
	_, err = io.Copy(toFd, fromFd)
	return err
}

// run runs a single command and returns an error if it does not succeed.
// os/exec should have this function, to be honest.
func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
