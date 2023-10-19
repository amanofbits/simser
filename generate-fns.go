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
