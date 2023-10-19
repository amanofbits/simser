package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"io"
)

type outputFile struct {
	fs         *token.FileSet
	file       *ast.File
	srcTypes   *types.Info
	importDecl *ast.GenDecl
	comments   []string
}

func newOutput(pkgName string, srcTypes *types.Info) (o *outputFile) {
	importDecl := &ast.GenDecl{
		Tok: token.IMPORT, Lparen: 1, Rparen: 1,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: "\"io\""}},
		},
	}
	return &outputFile{
		fs:         token.NewFileSet(),
		importDecl: importDecl,
		file: &ast.File{
			Name:    &ast.Ident{Name: pkgName},
			Imports: []*ast.ImportSpec{},
			Decls:   []ast.Decl{importDecl},
		},
		srcTypes: srcTypes,
	}
}

// func (o *outputFile) addImport(imp string) {
// 	imp = strings.Trim(imp, "\"")
// 	for _, s := range o.importDecl.Specs {
// 		spec := s.(*ast.ImportSpec)
// 		if strings.Trim(spec.Path.Value, "\"") == imp {
// 			return
// 		}
// 	}
// 	o.importDecl.Specs = append(o.importDecl.Specs,
// 		&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", imp)}},
// 	)
// }

func (o *outputFile) addFileComment(cmt string) {
	if o.comments == nil {
		o.comments = []string{}
	}
	o.comments = append(o.comments, cmt)
}

func (o *outputFile) appendFuncDecl(s *ast.FuncDecl) {
	o.file.Decls = append(o.file.Decls, s)
}

func (o outputFile) write(w io.Writer) (n int, err error) {

	o.file.Imports = make([]*ast.ImportSpec, len(o.importDecl.Specs))
	for i, s := range o.importDecl.Specs {
		o.file.Imports[i] = s.(*ast.ImportSpec)
	}

	buf := bytes.Buffer{}

	for _, com := range o.comments {
		buf.WriteString(com)
		buf.WriteString("\n\n")
	}
	if err := format.Node(&buf, o.fs, o.file); err != nil {
		return 0, fmt.Errorf("error formatting code, %w", err)
	}

	return w.Write(buf.Bytes())
}
