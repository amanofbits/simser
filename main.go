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
	"path/filepath"
	"strings"

	"github.com/am4n0w4r/simser/internal/generator"
	myParser "github.com/am4n0w4r/simser/internal/parser"
)

type config struct {
	targetFile  string
	rawTypes    string
	outputFile  string
	readFnName  string
	writeFnName string
}

func getConfig() (c config, err error) {

	c.targetFile = os.Getenv("GOFILE")
	c.targetFile, err = filepath.Abs(c.targetFile)
	if err != nil {
		return c, err
	}
	log.Printf("Target file: %s", c.targetFile)

	fi, err := os.Stat(c.targetFile)
	if err != nil {
		return c, fmt.Errorf("target file error, %w", err)
	}
	if !fi.Mode().IsRegular() {
		return c, fmt.Errorf("target file is not a regular file")
	}

	flag.StringVar(&c.rawTypes, "types", "", "comma-separated struct types to use")
	flag.StringVar(&c.outputFile, "output", "", "name of output file")
	flag.StringVar(&c.readFnName, "read-fn-name", "LoadFrom", "name of deserializing (read) function")
	flag.StringVar(&c.writeFnName, "write-fn-name", "SaveTo", "name of serializing (write) function")

	flag.Parse()

	if c.outputFile == "" {
		c.outputFile = fmt.Sprintf("%s.simser.g.go", strings.TrimSuffix(c.targetFile, ".go"))
	}

	return c, nil
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	acceptor, err := newTypeAcceptor(strings.Split(cfg.rawTypes, ","))
	if err != nil {
		log.Fatal(err)
	}

	file, err := myParser.Parse(cfg.targetFile)
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
		if err := generator.GenStructCode(s, output, cfg.readFnName, cfg.writeFnName); err != nil {
			log.Fatal(err)
		}
		log.Print("Done.")
	}

	if err := writeOutputFile(output, cfg.outputFile); err != nil {
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
	log.Print("Target file written")
	return nil
}
