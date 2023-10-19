package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

func analyzeInput(targetFile string) (*ast.File, *types.Info, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, targetFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var typesInfo = types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Uses:  make(map[*ast.Ident]types.Object),
		Defs:  make(map[*ast.Ident]types.Object),
		// Scopes: make(map[ast.Node]*types.Scope),
	}
	config := types.Config{Importer: importer.Default()}
	_, err = config.Check("", fs, []*ast.File{f}, &typesInfo)
	return f, &typesInfo, err
}

type structInfo struct {
	name   string
	decl   *ast.GenDecl
	spec   *ast.TypeSpec
	typ    *ast.StructType
	fields []fieldInfo
}

func selectTypes(f *ast.File, names map[string]any) ([]structInfo, error) {

	_, genAll := names["all"]
	genAll = genAll && len(names) == 1

	typs := []structInfo{}
	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			if decl.Tok != token.TYPE {
				continue
			}

			var spec *ast.TypeSpec
			var specTyp *ast.StructType
			for _, s := range decl.Specs {
				var ok bool
				if spec, ok = s.(*ast.TypeSpec); ok {
					t, ok := spec.Type.(*ast.StructType)
					if !ok {
						spec = nil
						continue
						// return typs, fmt.Errorf("type %s is not a struct type but %v", spec.Name.Name, t)
					}
					specTyp = t
					break
				}
				spec = nil
				specTyp = nil
			}

			if spec == nil {
				continue
			}

			typename := getTypeName(decl)
			if _, ok := names[typename]; !genAll && !ok {
				continue
			}
			delete(names, typename)

			typs = append(typs, structInfo{
				name:   typename,
				decl:   decl,
				spec:   spec,
				typ:    specTyp,
				fields: []fieldInfo{},
			})
		default:
			continue
		}
	}
	if !genAll && len(names) > 0 {
		errs := []error{}
		for k := range names {
			errs = append(errs, fmt.Errorf("type %s were not found in file", k))
		}
		return typs, errors.Join(errs...)
	}
	return typs, nil
}

func getTypeName(decl *ast.GenDecl) string {
	for _, s := range decl.Specs {
		if typeSpec, ok := s.(*ast.TypeSpec); ok {
			return typeSpec.Name.Name
		}
	}
	return ""
}

type fieldInfo struct {
	name         string
	size         int
	typeToRead   *types.Basic
	typeToAssign types.Type
}

func getFieldInfo(f *ast.Field, typesInfo *types.Info) (info fieldInfo, err error) {
	t := typesInfo.TypeOf(f.Type)
	loop := true
	for loop {
		switch typ := t.(type) {
		case *types.Basic:
			loop = false
			info.typeToRead = typ
			if info.typeToAssign == nil {
				info.typeToAssign = typ
			}
			continue
		case *types.Named:
			info.typeToAssign = typ
			t = typ.Underlying()
		case *types.Slice: // how to know how many bytes to read?
			return info, fmt.Errorf("unsupported type %T, %v", t, t)
		case *types.Array:
			return info, fmt.Errorf("not supported yet type %T, %v", t, t)
		default:
			return info, fmt.Errorf("unsupported type %T, %v", t, t)
		}
	}

	typ := t.(*types.Basic)
	switch typ.Kind() {
	case types.Int8, types.Uint8:
		info.size = 1
	case types.Int16, types.Uint16:
		info.size = 2
	case types.Int32, types.Uint32, types.Float32:
		info.size = 4
	case types.Int64, types.Uint64, types.Float64:
		info.size = 8
	case types.Int, types.Uint:
		return info, fmt.Errorf("%v: type size is platform-dependent, please choose exact-sized equivalent", t)
	default:
		return info, fmt.Errorf("type not found, %v", t)
	}
	return info, nil
}
