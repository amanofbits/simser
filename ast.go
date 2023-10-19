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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func astMakeByteBuffer(name string, length int, cap int) *ast.AssignStmt {
	if cap < length {
		cap = length
	}
	return &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(name)},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					ast.NewIdent("[]byte"),
					ast.NewIdent(strconv.Itoa(length)),
					ast.NewIdent(strconv.Itoa(cap)),
				},
			},
		},
	}
}

func astFromStatements(src string) ([]ast.Stmt, error) {
	node, err := parser.ParseExpr(fmt.Sprintf("func (r io.Reader) { %s }", src))
	if err != nil {
		return nil, err
	}
	fn := node.(*ast.FuncLit)
	return fn.Body.List, nil
}

func asIdents(names ...string) []*ast.Ident {
	idents := []*ast.Ident{}
	for _, n := range names {
		idents = append(idents, ast.NewIdent(n))
	}
	return idents
}

func funcDeclRead(typ structInfo) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("ReadFrom"),
		Recv: &ast.FieldList{List: []*ast.Field{{Names: asIdents("o"), Type: &ast.StarExpr{X: ast.NewIdent(typ.name)}}}},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "r"}},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("io"),
							Sel: ast.NewIdent("Reader"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "n"}},
						Type:  ast.NewIdent("int64"),
					},
					{
						Names: []*ast.Ident{{Name: "err"}},
						Type:  ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			Lbrace: 1, Rbrace: 1,
			List: []ast.Stmt{},
		},
	}
}

func funcDeclWrite(typ structInfo) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("WriteTo"),
		Recv: &ast.FieldList{List: []*ast.Field{{Names: asIdents("s"), Type: &ast.StarExpr{X: ast.NewIdent(typ.name)}}}},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "w"}},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("io"),
							Sel: ast.NewIdent("Writer"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "n"}},
						Type:  ast.NewIdent("int64"),
					},
					{
						Names: []*ast.Ident{{Name: "err"}},
						Type:  ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			Lbrace: 1, Rbrace: 1,
			List: []ast.Stmt{},
		},
	}
}
