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
	"go/ast"
	"log"
)

func generateFunctions(typs []structInfo, output *outputFile) error {

	for _, typ := range typs {
		log.Printf("Processing %s...", typ.name)
		if err := generateForStruct(typ, output); err != nil {
			return err
		}
		log.Print("Done.")
	}

	return nil
}

func generateForStruct(typ structInfo, output *outputFile) error {
	if len(typ.fields) == 0 {
		return nil
	}
	// Gen func typenameRead(io.Reader) (*typename, error)
	{
		maxBufLen := 1
		for _, f := range typ.fields {
			if maxBufLen < f.size {
				maxBufLen = f.size
			}
		}
		declReadFn := funcDeclRead(typ)
		if len(typ.fields) != 0 {

			declReadFn.Body.List = append(declReadFn.Body.List,
				astMakeByteBuffer("b", maxBufLen, maxBufLen))

			for _, field := range typ.fields {
				// if field.size > 1 {
				// 	output.addImport("encoding/binary")
				// }
				temp, err := tplReadField(field, "b")
				if err != nil {
					return err
				}
				stmt, err := astFromStatements(temp)
				if err != nil {
					return err
				}
				bstmt := ast.BlockStmt{Lbrace: 1, Rbrace: 1, List: stmt}
				declReadFn.Body.List = append(declReadFn.Body.List, &bstmt)
			}
		}

		declReadFn.Body.List = append(declReadFn.Body.List, &ast.ReturnStmt{
			Results: []ast.Expr{ast.NewIdent("n"), ast.NewIdent("err")},
		})
		output.appendFuncDecl(declReadFn)
	}

	// Gen func typenameWrite(io.Writer) (n int, err error)
	{
		totalSize := 0
		for _, f := range typ.fields {
			totalSize += f.size
		}
		declWriteFn := funcDeclWrite(typ)
		if len(typ.fields) != 0 {

			declWriteFn.Body.List = append(declWriteFn.Body.List,
				astMakeByteBuffer("b", 0, totalSize))

			for _, field := range typ.fields {
				// if field.size > 1 {
				// 	output.addImport("encoding/binary")
				// }
				temp, err := tplWriteField(field, "b")
				if err != nil {
					return err
				}
				stmt, err := astFromStatements(temp)
				if err != nil {
					return err
				}
				bstmt := ast.BlockStmt{List: stmt}
				declWriteFn.Body.List = append(declWriteFn.Body.List, &bstmt)
			}
		}

		stmt, err := astFromStatements("nW, err := w.Write(b);")
		if err != nil {
			return err
		}
		declWriteFn.Body.List = append(declWriteFn.Body.List, stmt...)

		declWriteFn.Body.List = append(declWriteFn.Body.List, &ast.ReturnStmt{
			Results: []ast.Expr{ast.NewIdent("int64(nW)"), ast.NewIdent("err")},
		})
		output.appendFuncDecl(declWriteFn)
	}

	return nil
}
