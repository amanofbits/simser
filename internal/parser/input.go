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
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"

	"github.com/am4n0w4r/simser/internal/domain"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

func Parse(targetFile string) (file *InputFile, err error) {
	file = &InputFile{
		Path: targetFile,
	}

	file.Modfile, file.Moddir, err = GetClosestModFile(filepath.Dir(targetFile))
	if err != nil {
		return file, err
	}

	importPath, err := pathToImport(targetFile, file.Moddir, file.Modfile.Module.Mod.Path)
	if err != nil {
		return file, err
	}

	pfset := token.NewFileSet()
	pkgs, err := loadPackage(importPath, file.Moddir, pfset)
	if err != nil {
		return file, err
	}

	switch len(pkgs) {
	case 0:
		return file, errors.New("found 0 file packages")
	case 1:
		file.Pkg = pkgs[0]
	default:
		return file, fmt.Errorf("%d ambiguous packages found for file", len(pkgs))
	}
	file.SyntaxIdx = slices.Index(file.Pkg.CompiledGoFiles, file.Path) // let's hope the order is the same
	if file.SyntaxIdx < 0 {
		return file, errors.New("file not found in package. Looks like programming error")
	}

	return file, nil
}

type InputFile struct {
	Path      string            // Path of the file
	Modfile   *modfile.File     // File's module
	Moddir    string            // Directory where module file is present
	Pkg       *packages.Package // Parsed file's package
	SyntaxIdx int               // index of file ast in pkg.Syntax
}

func (fi InputFile) Ast() *ast.File { return fi.Pkg.Syntax[fi.SyntaxIdx] }

type TypeAcceptor interface {
	// Should return true if type is required to be serializable
	Accepts(tspec *ast.TypeSpec) bool
	// Should return true if no more types are expected
	IsDrained() bool
	// Should return true if acceptor has no type restrictions
	AcceptsAll() bool
	fmt.Stringer
}

func (f InputFile) GetInputStructs(acceptor TypeAcceptor) ([]domain.InputStruct, error) {
	filtered, err := f.filterInputStructs(acceptor)
	if err != nil {
		return nil, err
	}
	if !acceptor.IsDrained() && !acceptor.AcceptsAll() {
		return nil, fmt.Errorf("types '%s' were requested but not found", acceptor)
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no types %s found", acceptor)
	}

	return analyzeStructs(filtered, f)
}

type filteredStruct struct {
	name     string
	astType  *ast.StructType
	typeInfo *types.Struct
}

func (s filteredStruct) fieldCount() int { return s.typeInfo.NumFields() }

func (f InputFile) filterInputStructs(acceptor TypeAcceptor) (structs []filteredStruct, err error) {
	fileAst := f.Ast()
	structs = []filteredStruct{}

	for _, d := range fileAst.Decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			continue
		}
		for _, s := range decl.Specs {
			tspec, ok := s.(*ast.TypeSpec)
			if !ok || !acceptor.Accepts(tspec) {
				continue
			}
			structName := tspec.Name.Name
			astStruct, ok := tspec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			ti := f.Pkg.TypesInfo.TypeOf(astStruct)
			if ti == nil {
				return structs, fmt.Errorf("cannot get type info for struct %s", structName)
			}
			structType, ok := ti.(*types.Struct)
			if !ok {
				return structs, fmt.Errorf("type %v is not a struct but %T", ti, ti)
			}

			structs = append(structs, filteredStruct{
				name:     structName,
				astType:  astStruct,
				typeInfo: structType,
			})
			if acceptor.IsDrained() {
				break
			}
		}
		if acceptor.IsDrained() {
			break
		}
	}
	return structs, nil
}

func pathToImport(file, moddir, modname string) (string, error) {
	path, err := filepath.Rel(moddir, file)
	if err != nil {
		return "", err
	}
	path = filepath.ToSlash(filepath.Join(modname, path))
	return filepath.Dir(path), nil
}
