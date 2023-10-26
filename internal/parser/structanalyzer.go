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
	"go/types"
	"strings"

	"github.com/am4n0w4r/simser/internal/domain"
)

func analyzeStructs(filtered []filteredStruct, f InputFile) (structs []domain.InputStruct, err error) {

	structs = make([]domain.InputStruct, len(filtered))
	for i, fs := range filtered {
		s, err := analyzeStruct(fs, f.Pkg.PkgPath)
		if err != nil {
			return nil, err
		}
		structs[i] = *s
	}

	return structs, nil
}

func analyzeStruct(fs filteredStruct, pkgPath string) (s *domain.InputStruct, err error) {

	s = domain.NewInputStruct(fs.name, fs.typeInfo)
	fields := make([]domain.StructField, fs.fieldCount())

	for i := 0; i < fs.fieldCount(); i++ {
		sField := fs.typeInfo.Field(i)
		tag := *newStructTag()
		if err := tag.parse(fs.astType.Fields.List[i].Tag); err != nil {
			return nil, err
		}

		fTyp, err := getFieldType(sField.Type(), pkgPath, &tag)

		field := domain.NewStructField(sField.Name(), fTyp, tag.values)
		if err != nil {
			return nil, fmt.Errorf("field '%s.%s %s': %w", s.Name(), sField.Name(), sField.Type(), err)
		}
		fields[i] = field
	}
	s.SetFields(fields)
	return s, nil
}

func getFieldType(t types.Type, trimPkgPath string, tag *structTag) (ft domain.FieldType, err error) {
	switch typ := t.(type) {

	case *types.Basic:
		name, err := getTypeName(typ, trimPkgPath)
		if err != nil {
			return nil, err
		}
		size, err := getTypeSize(typ)
		if err != nil {
			return nil, err
		}
		return domain.NewSimpleFieldType(name, size, nil), nil

	case *types.Named:
		name, err := getTypeName(typ, trimPkgPath)
		if err != nil {
			return nil, err
		}
		size, err := getTypeSize(typ)
		if err != nil {
			return nil, err
		}

		var underlying types.Type = typ
		for underlying.Underlying() != underlying {
			underlying = typ.Underlying()
		}
		ut, ok := underlying.(*types.Basic)
		if !ok {
			return nil, errors.Join(domain.ErrUnsupportedType, fmt.Errorf("underlying type is not a basic type, but %T", underlying))
		}
		return domain.NewSimpleFieldType(
			name,
			size,
			domain.NewSimpleFieldType(ut.Name(), size, nil),
		), nil

	case *types.Array:
		if typ.Len() < 0 {
			return nil, errors.Join(domain.ErrUnsupportedType, fmt.Errorf("array type with unknown length, %T, %v", typ, typ))
		}
		el, err := getFieldType(typ.Elem(), trimPkgPath, tag)
		if err != nil {
			return nil, fmt.Errorf("failed to get array element type, %w", err)
		}
		return domain.NewArrayFieldType(int(typ.Len()), el), nil

	case *types.Slice:
		lenTag, ok, err := tag.getLenExpr()
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.Join(domain.ErrUnsupportedType,
				errors.New("slice length cannot be determined. use 'len' tag attribute to supply an expression (`simser:\"len=o.fieldName\"`)"))
		}
		if lenTag == "" {
			return nil, errors.Join(domain.ErrUnsupportedType, errors.New("empty length expression"))
		}

		if err != nil {
			return nil, errors.Join(domain.ErrUnsupportedType, err)
		}

		el, err := getFieldType(typ.Elem(), trimPkgPath, tag)
		if err != nil {
			return nil, errors.Join(domain.ErrUnsupportedType, fmt.Errorf("failed to get slice element type, %T", typ.Elem()))
		}
		return domain.NewSliceFieldType(lenTag, el), nil

	default:
		return nil, errors.Join(domain.ErrUnsupportedType, fmt.Errorf("%T, %v", typ, typ))
	}
}

func getTypeName(t types.Type, trimPkgPath string) (name string, err error) {
	trimPkgPath = trimPkgPath + "."

	switch typ := t.(type) {
	case *types.Basic:
		name = typ.Name()
		name = strings.TrimPrefix(name, trimPkgPath)
		return name, nil

	case *types.Named:
		name = typ.String()
		name = strings.TrimPrefix(name, trimPkgPath)
		return name, nil

	default:
		return "", fmt.Errorf("failed to get name and size for %T", typ)
	}
}

func getTypeSize(t types.Type) (int, error) {
	tOrig := t
	for t != t.Underlying() {
		t = t.Underlying()
	}
	bt, ok := t.(*types.Basic)
	if !ok {
		return 0, fmt.Errorf("type %v has no underlying basic type, but %T", tOrig, t)
	}
	switch bt.Kind() {
	case types.Int8, types.Uint8:
		return 1, nil
	case types.Int16, types.Uint16:
		return 2, nil
	case types.Int32, types.Uint32, types.Float32:
		return 4, nil
	case types.Int64, types.Uint64, types.Float64:
		return 8, nil
	case types.Int, types.Uint:
		return 0, fmt.Errorf("%v: type size differs based on platform, please choose exact-sized one", t)
	default:
		return 0, fmt.Errorf("type not found, %v", t)
	}
}
