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

// Domain model, outputs of parser, inputs for generator.
// No dependencies here.
package domain

import (
	"errors"
	"fmt"
	"go/types"
	"strings"
)

var ErrUnsupportedType = errors.New("unsupported type")

type InputStruct struct {
	name   string
	typ    *types.Struct
	fields []StructField
}

func NewInputStruct(name string, typ *types.Struct) InputStruct {
	return InputStruct{
		name: name,
		typ:  typ,
	}
}

func (s InputStruct) Name() string        { return s.name }
func (s InputStruct) Type() *types.Struct { return s.typ }

func (s InputStruct) Field(idx int) StructField  { return s.fields[idx] }
func (s InputStruct) FieldCount() int            { return len(s.fields) }
func (s *InputStruct) SetFields(f []StructField) { s.fields = f }

//

type StructField struct {
	name string
	typ  FieldType
}

func NewStructField(name string, typ FieldType) StructField {
	return StructField{
		name: name,
		typ:  typ,
	}
}

func (f StructField) Name() string    { return f.name }
func (f StructField) Type() FieldType { return f.typ }

//

type FieldType interface {
	Name() string
	Size() int // Total size of the field
	// Don't confuse with IsInt. Checks if type is intX OR uintX
	IsInteger() bool
}

type SimpleFieldType struct {
	name       string
	size       int
	underlying *SimpleFieldType
}

func NewSimpleFieldType(name string, size int, underlying *SimpleFieldType) *SimpleFieldType {
	return &SimpleFieldType{
		name:       name,
		size:       size,
		underlying: underlying,
	}
}

func (t SimpleFieldType) Name() string                 { return t.name }
func (t SimpleFieldType) Size() int                    { return t.size }
func (t SimpleFieldType) BitSize() int                 { return t.size * 8 }
func (t SimpleFieldType) Underlying() *SimpleFieldType { return t.underlying }

func (bt SimpleFieldType) IsInteger() bool {
	tmp := &bt
	for tmp.underlying != nil {
		tmp = tmp.underlying
	}
	return strings.HasPrefix(tmp.name, "int") || strings.HasPrefix(tmp.name, "uint") || tmp.name == "byte"
}

type ArrayFieldType struct {
	length int
	elType FieldType
}

func NewArrayFieldType(length int, el FieldType) *ArrayFieldType {
	return &ArrayFieldType{
		length: length,
		elType: el,
	}
}

func (at ArrayFieldType) Name() string      { return fmt.Sprintf("[]%s", at.elType.Name()) }
func (at ArrayFieldType) Size() int         { return at.elType.Size() * at.length }
func (at ArrayFieldType) Length() int       { return at.length }
func (at ArrayFieldType) ElType() FieldType { return at.elType }
func (at ArrayFieldType) IsInteger() bool   { return false }
