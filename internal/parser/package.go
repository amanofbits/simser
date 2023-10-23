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
	"fmt"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

func loadPackage(importName, moddir string, fset *token.FileSet) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedCompiledGoFiles | packages.NeedDeps | packages.NeedFiles | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedTypes | packages.NeedName,
		Dir:  moddir,
		Fset: fset,
	}

	importName = strings.Trim(importName, "\"")
	pkgs, err := packages.Load(cfg, importName)
	if err != nil {
		return nil, fmt.Errorf("error loading package %s: %w", importName, err)
	}
	return pkgs, nil
}
