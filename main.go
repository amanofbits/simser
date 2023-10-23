// Copyright 2023 am4n0w4r
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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/am4n0w4r/simser/internal/generator"
	myParser "github.com/am4n0w4r/simser/internal/parser"
)

type flags struct {
	rawTypes    string
	outputFile  string
	readFnName  string
	writeFnName string
}

func getFlags() (f flags) {

	flag.StringVar(&f.rawTypes, "types", "", "comma-separated struct types to use")
	flag.StringVar(&f.outputFile, "output", "", "name of output file")
	flag.StringVar(&f.readFnName, "read-fn-name", "LoadFrom", "name of deserializing (read) function")
	flag.StringVar(&f.writeFnName, "write-fn-name", "SaveTo", "name of serializing (write) function")

	flag.Parse()

	return f
}

func main() {
	targetFile := os.Getenv("GOFILE")

	f := getFlags()
	if f.outputFile == "" {
		f.outputFile = fmt.Sprintf("%s.simser.g.go", strings.TrimSuffix(targetFile, ".go"))
	}

	acceptor, err := newTypeAcceptor(strings.Split(f.rawTypes, ","))
	if err != nil {
		log.Fatal(err)
	}

	file, err := myParser.Parse(targetFile)
	if err != nil {
		log.Fatal(err)
	}

	inputStructs, err := file.GetInputStructs(acceptor)
	if err != nil {
		log.Fatal(err)
	}

	output := generator.NewOutput(file.Pkg)

	for _, s := range inputStructs {
		log.Printf("Processing %s...", s.Name())
		if err := generator.GenStructCode(s, output, f.readFnName, f.writeFnName); err != nil {
			log.Fatal(err)
		}
		log.Print("Done.")
	}

	if err := writeOutputFile(output, f.outputFile); err != nil {
		log.Fatal(err)
	}
}

func writeOutputFile(output *generator.Output, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file, %w", err)
	}
	defer file.Close()

	if _, err := output.WriteTo(file); err != nil {
		return fmt.Errorf("failed to write outpput file, %w", err)
	}
	log.Printf("file %s written", filename)
	return nil
}
