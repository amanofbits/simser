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
	"strconv"
	"strings"
)

var ErrUnsupportedType = errors.New("unsupported type")

type InputStruct struct {
	name   string
	typ    *types.Struct
	fields []StructField
}

func NewInputStruct(name string, typ *types.Struct) *InputStruct {
	return &InputStruct{
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
	tag  map[string]string
}

func NewStructField(name string, typ FieldType, tag map[string]string) StructField {
	return StructField{
		name: name,
		typ:  typ,
		tag:  tag,
	}
}

func (f StructField) Name() string    { return f.name }
func (f StructField) Type() FieldType { return f.typ }

// Type

type FieldType interface {
	Name() string
	// Total size of the field.
	// If Size() < 0, then size is not fixed, and SizeExpr() should be used.
	Size() int
	// Expression that determines the size of the field.
	// Should always be constructed in a way to return int type.
	SizeExpr() string
	// Don't confuse with IsInt. Checks if type is intX OR uintX OR byte
	IsInteger() bool
	// Should return true if the type is an array or slice
	IsSequence() bool
}

// Field type that is a sequence, i.e. array or slice, which have a number of certain elements
type SequenceFieldType interface { // No need to make this a more generic collection, as [de]serialization is sequential
	ElType() FieldType
	LenExpr() string
}

// Simple

type SimpleFieldType struct {
	name       string
	size       int
	underlying *SimpleFieldType
}

func NewSimpleFieldType(name string, size int, underlying *SimpleFieldType) *SimpleFieldType {
	if size < 1 {
		panic(fmt.Sprintf("simple type size < 1. This should be caught earlier. Please, file a bug to the repo. Type %s, size %d",
			name, size))
	}
	return &SimpleFieldType{
		name:       name,
		size:       size,
		underlying: underlying,
	}
}

func (t SimpleFieldType) Name() string { return t.name }

func (t SimpleFieldType) Size() int                    { return t.size }
func (t SimpleFieldType) SizeExpr() string             { return strconv.Itoa(t.size) }
func (t SimpleFieldType) BitSize() int                 { return t.size * 8 }
func (t SimpleFieldType) Underlying() *SimpleFieldType { return t.underlying }
func (t SimpleFieldType) IsSequence() bool             { return false }

func (bt SimpleFieldType) IsInteger() bool {
	tmp := &bt
	for tmp.underlying != nil {
		tmp = tmp.underlying
	}
	return strings.HasPrefix(tmp.name, "int") || strings.HasPrefix(tmp.name, "uint") || tmp.name == "byte"
}

// Array

type ArrayFieldType struct {
	length int
	elType FieldType
}

func NewArrayFieldType(length int, el FieldType) *ArrayFieldType {
	if length < 0 {
		panic(fmt.Sprintf("array length < 1. This should be caught earlier. Please, file a bug to the repo. Length %d",
			length))
	}
	return &ArrayFieldType{
		length: length,
		elType: el,
	}
}

func (t ArrayFieldType) Name() string { return fmt.Sprintf("[]%s", t.elType.Name()) }

func (t ArrayFieldType) Size() int {
	if !IsFixedSize(t.elType.Size()) {
		return -1
	}
	return t.elType.Size() * t.length
}
func (t ArrayFieldType) SizeExpr() string {
	if IsFixedSize(t) {
		return strconv.Itoa(t.elType.Size() * t.length)
	}
	return fmt.Sprintf("(%s) * %d", t.elType.SizeExpr(), t.length)
}
func (t ArrayFieldType) Length() int       { return t.length }
func (t ArrayFieldType) LenExpr() string   { return strconv.Itoa(t.length) }
func (t ArrayFieldType) ElType() FieldType { return t.elType }
func (t ArrayFieldType) IsInteger() bool   { return false }
func (t ArrayFieldType) IsSequence() bool  { return true }

// Slice

type SliceFieldType struct {
	lenExpr string // expression used to calculate length.
	elType  FieldType
}

func NewSliceFieldType(lenExpr string, el FieldType) *SliceFieldType {
	return &SliceFieldType{
		lenExpr: strings.TrimSpace(lenExpr),
		elType:  el,
	}
}

func (t SliceFieldType) Name() string { return fmt.Sprintf("[]%s", t.elType.Name()) }

func (t SliceFieldType) Size() int { return -1 }
func (t SliceFieldType) SizeExpr() string {
	return fmt.Sprintf("%s * %s", ParenthesizeIntExpr(t.elType.SizeExpr()), ParenthesizeIntExpr(t.lenExpr))
}
func (t SliceFieldType) LenExpr() string   { return t.lenExpr }
func (t SliceFieldType) ElType() FieldType { return t.elType }
func (t SliceFieldType) IsInteger() bool   { return false }
func (t SliceFieldType) IsSequence() bool  { return true }
