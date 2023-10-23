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

package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func GetClosestModFile(startPath string) (file *modfile.File, moddir string, err error) {
	fInfo, err := os.Stat(startPath)
	if err != nil {
		return nil, startPath, err
	}
	if !fInfo.IsDir() {
		startPath = filepath.Dir(startPath)
	}

	moddir, err = findDirWithFile("go.mod", startPath)
	if err != nil {
		return nil, moddir, fmt.Errorf("failed to find go.mod starting with %s", startPath)
	}

	raw, err := os.ReadFile(filepath.Join(moddir, "go.mod"))
	if err != nil {
		return nil, moddir, fmt.Errorf("failed to read go.mod, %w", err)
	}

	modFile, err := modfile.Parse("go.mod", raw, nil)
	if err != nil {
		return modFile, moddir, fmt.Errorf("failed to parse module, %w", err)
	}
	return modFile, moddir, nil
}

func findDirWithFile(filename, startDir string) (dir string, err error) {
	dir = startDir

	_, err = os.Stat(filepath.Join(dir, filename))
	if err == nil {
		return dir, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return dir, err
	}

	for {
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return dir, os.ErrNotExist
		}
		dir = parentDir

		_, err := os.Stat(filepath.Join(dir, filename))
		if err == nil {
			return dir, nil
		}
	}
}
