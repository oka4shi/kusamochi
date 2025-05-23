//go:build ignore

// The original code was written by Go authors.
// This source code is governed by a BSD-style license:
/*
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const suffix = ".ttf"

func main() {
	ttfs, err := os.Open("ttfs")
	if err != nil {
		log.Panicln(err)
	}
	defer ttfs.Close()

	infos, err := ttfs.Readdir(-1)
	if err != nil {
		log.Panicln(err)
	}
	for _, info := range infos {
		ttfName := info.Name()
		if !strings.HasSuffix(ttfName, suffix) {
			continue
		}
		do(ttfName)
	}
}

func do(ttfName string) {
	fontName := fontName(ttfName)
	pkgName := pkgName(ttfName)
	if err := os.Mkdir(pkgName, 0777); err != nil && !os.IsExist(err) {
		log.Panicln(err)
	}
	src, err := os.ReadFile(filepath.Join("ttfs", ttfName))
	if err != nil {
		log.Panicln(err)
	}

	oflFile := oflFile(ttfName)
	ofl, err := os.ReadFile(filepath.Join("ttfs", oflFile))
	if err != nil {
		log.Panicln(err)
	}

	b := new(bytes.Buffer)
	fmt.Fprintf(b, "// generated by go run gen.go; DO NOT EDIT\n\n")
	fmt.Fprintf(b, "// Package %s provides the %q TrueType font.\n", pkgName, fontName)
	fmt.Fprintf(b, "\n/*\n")
	b.Write(ofl)
	fmt.Fprintf(b, "*/\n\n")
	fmt.Fprintf(b, "package %s\n\n", pkgName)
	fmt.Fprintf(b, "// TTF is the data for the %q TrueType font.\n", fontName)
	fmt.Fprintf(b, "var TTF = []byte{")
	for i, x := range src {
		if i&15 == 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(b, "%#02x,", x)
	}
	fmt.Fprintf(b, "\n}\n")

	dst, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgName, "data.go"), dst, 0666); err != nil {
		log.Fatal(err)
	}
}

// fontName maps "Go-Regular.ttf" to "Go Regular".
func fontName(ttfName string) string {
	s := ttfName[:len(ttfName)-len(suffix)]
	s = strings.Replace(s, "-", " ", -1)
	return s
}

// pkgName maps "Go-Regular.ttf" to "goregular".
func pkgName(ttfName string) string {
	s := ttfName[:len(ttfName)-len(suffix)]
	s = strings.Replace(s, "-", "", -1)
	s = strings.ToLower(s)
	return s
}

func oflFile(ttfName string) string {
	return "OFL_" + strings.Replace(ttfName, ".ttf", ".txt", -1)
}
