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

func getFields(s domain.InputStruct, pkgPath string) (fields []domain.StructField, err error) {
	fields = make([]domain.StructField, s.Type().NumFields())

	for i := 0; i < s.Type().NumFields(); i++ {
		sField := s.Type().Field(i)
		fTyp, err := getFieldType(sField.Type(), pkgPath)
		field := domain.NewStructField(sField.Name(), fTyp)
		if err != nil {
			return nil, fmt.Errorf("field '%s.%s %s': %w", s.Name(), sField.Name(), sField.Type(), err)
		}
		fields[i] = field
	}
	return fields, nil
}

func getFieldType(t types.Type, trimPkgPath string) (ft domain.FieldType, err error) {
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
		ub, ok := underlying.(*types.Basic)
		if !ok {
			return nil, errors.Join(domain.ErrUnsupportedType, fmt.Errorf("underlying type is not a basic type, but %T", underlying))
		}
		return domain.NewSimpleFieldType(
			name,
			size,
			domain.NewSimpleFieldType(ub.Name(), size, nil),
		), nil

	case *types.Array:
		if typ.Len() < 0 {
			return nil, errors.Join(domain.ErrUnsupportedType, errors.New("slice length cannot be determined"))
		}
		el, err := getFieldType(typ.Elem(), trimPkgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get array element type, %T", typ.Elem())
		}
		return domain.NewArrayFieldType(int(typ.Len()), el), nil

	case *types.Slice:
		return nil, errors.Join(domain.ErrUnsupportedType, errors.New("slice length cannot be determined"))

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
