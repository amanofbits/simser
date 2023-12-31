// Copyright 2023 amanofbits
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"fmt"
	"go/parser"
	"go/printer"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Output struct {
	pkg     *packages.Package
	header  fstringBuilder
	imports map[string]string
	code    fstringBuilder
}

func NewOutput(pkg *packages.Package) (o *Output) {
	o = &Output{
		pkg:     pkg,
		header:  fstringBuilder{},
		imports: map[string]string{},
		code:    fstringBuilder{},
	}
	o.header.WriteFString("// Code generated by \"%s %s\"; DO NOT EDIT.\n\n", filepath.Base(os.Args[0]), strings.Join(os.Args[1:], " "))
	o.header.WriteFString("package %s\n", pkg.Name)

	return o
}

func (o *Output) AppendImport(imp string) {
	imp = strings.Trim(imp, "\"")
	o.imports[imp] = ""
}

func (o *Output) AppendNamedImport(name, imp string) {
	imp = strings.Trim(imp, "\"")
	o.imports[imp] = name
}

func (o *Output) AppendF(text string, args ...any) *Output {
	o.code.WriteFString(text, args...)
	return o
}
func (o *Output) Append(text string) *Output {
	o.code.WriteString(text)
	return o
}

func (o *Output) LF() *Output {
	o.code.WriteByte('\n')
	return o
}

func (o Output) WriteTo(w io.Writer) (n int64, err error) {
	src := fstringBuilder{}

	src.WriteString(o.header.String())

	src.WriteString("import (\n")
	for imp, name := range o.imports {
		src.WriteFString("%s \"%s\"\n", name, imp)
	}
	src.WriteString(")\n")

	src.WriteString(o.code.String())

	srcStr := src.String()

	// Parse src to check for errors
	f, err := parser.ParseFile(o.pkg.Fset, "", srcStr, parser.ParseComments)
	if err != nil {
		return 0, err
	}

	// Format src and print to a writer
	if err := printer.Fprint(w, o.pkg.Fset, f); err != nil {
		return 0, fmt.Errorf("error formatting code, %w", err)
	}

	return 0, nil
}
