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

func astDefineVarDefault(name, typ string) ast.Stmt {
	stmt, err := astFromStatements(fmt.Sprintf("var %s %s", name, typ))
	if err != nil {
		panic(err)
	}
	return stmt[0]
}

func astReadToSlice(targetSliceName string, targetSliceLen int) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("n"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("Read"),
					},
					Args: []ast.Expr{ast.NewIdent(targetSliceName)},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("err"),
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{ast.NewIdent("nil"), ast.NewIdent("err")},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("n"),
				Op: token.NEQ,
				Y:  ast.NewIdent(strconv.Itoa(targetSliceLen)),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{ast.NewIdent("nil"), ast.NewIdent("err")},
					},
				},
			},
		},
	}
}

func astRetReadValue(sliceName string) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("val")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{&ast.IndexExpr{
				X:     ast.NewIdent(sliceName),
				Index: ast.NewIdent("0"),
			}},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{&ast.UnaryExpr{
				Op: token.AND,
				X:  ast.NewIdent("val"),
			}, ast.NewIdent("nil")},
		},
	}
}

func astFromStatements(src string) ([]ast.Stmt, error) {
	// fset := token.NewFileSet()
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
