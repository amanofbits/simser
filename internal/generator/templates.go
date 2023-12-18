// Copyright 2023 amanofbits
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

package generator

import (
	"errors"
	"fmt"

	"github.com/amanofbits/simser/internal/domain"
)

// read toRead bytes
func tpl_ReadBytesIntoBuf(bufName string) string {
	return fmt.Sprintf(
		`nRead, err = io.ReadFull(r, %s[:toRead])
n += nRead
if err != nil {
	return n, err
}`, bufName)
}

func tpl_WriteField(f domain.StructField, bufName string) (t string, err error) {
	sb := fstringBuilder{}

	switch fType := f.Type().(type) {
	case *domain.SimpleFieldType:
		tpl_AppendSimpleTypeToBytes(bufName, f.Name(), fType, &sb)

	case *domain.ArrayFieldType:
		sb.WriteFString("for i:=0;i<len(o.%s);i++ {\n", f.Name())
		elType, ok := fType.ElType().(*domain.SimpleFieldType)
		if !ok {
			return sb.String(), errors.Join(domain.ErrUnsupportedType, errors.New("array of arrays are not supported"))
		}
		tpl_AppendSimpleTypeToBytes("b", fmt.Sprintf("%s[i]", f.Name()), elType, &sb)
		sb.WriteString("\n}")

	case *domain.SliceFieldType:
		sb.WriteFString("for i:=0;i<len(o.%s);i++ {\n", f.Name())
		elType, ok := fType.ElType().(*domain.SimpleFieldType)
		if !ok {
			return sb.String(), errors.Join(domain.ErrUnsupportedType, errors.New("slice of arrays are not supported"))
		}
		tpl_AppendSimpleTypeToBytes("b", fmt.Sprintf("%s[i]", f.Name()), elType, &sb)
		sb.WriteString("\n}")
	}

	return sb.String(), nil
}

func tpl_AppendSimpleTypeToBytes(bufName string, fieldName string, t *domain.SimpleFieldType, dst *fstringBuilder) {
	dst.WriteFString("%s = append(%s, ", bufName, bufName)

	intStr := fmt.Sprintf("uint%d", t.BitSize())
	for i := 0; i < t.Size(); i++ {
		dst.WriteString("byte(")
		if !t.IsInteger() {
			dst.WriteFString("%s(", intStr)
		}
		dst.WriteFString("o.%s", fieldName)
		if !t.IsInteger() {
			dst.WriteString(")")
		}
		if i != 0 {
			dst.WriteFString(">>%d", i*8)
		}
		dst.WriteString(")")
		if i != t.Size()-1 {
			dst.WriteString(",")
		}
	}
	dst.WriteString(")")
}

func tpl_ReadField(f domain.StructField, bufName string) (s string, err error) {
	sb := fstringBuilder{}

	switch fType := f.Type().(type) {

	case *domain.SimpleFieldType:
		sb.WriteFString("o.%s = ", f.Name())
		tpl_BytesToSimpleType("b", fType, &sb)

	case *domain.ArrayFieldType:
		elType, ok := fType.ElType().(*domain.SimpleFieldType)
		if !ok {
			return sb.String(), errors.Join(domain.ErrUnsupportedType, errors.New("array of arrays are not supported"))
		}

		sb.WriteFString("o.%s = [%d]%s{}\n", f.Name(), fType.Length(), elType.Name())

		sb.WriteFString("for i:=0;i<len(o.%s);i++ {\n", f.Name())
		sb.WriteFString("o.%s[i] = ", f.Name())
		tpl_BytesToSimpleType("b", elType, &sb)
		sb.WriteString("\n}")

	case *domain.SliceFieldType:
		elType, ok := fType.ElType().(*domain.SimpleFieldType)
		if !ok {
			return sb.String(), errors.Join(domain.ErrUnsupportedType, errors.New("array of arrays are not supported"))
		}
		sb.WriteFString("o.%s = make([]%s, sLen)\n", f.Name(), fType.ElType().Name())

		sb.WriteFString("for i:=0;i<len(o.%s);i++ {\n", f.Name())
		sb.WriteFString("o.%s[i] = ", f.Name())
		tpl_BytesToSimpleType("b", elType, &sb)
		sb.WriteString("\n}")

	default:
		return sb.String(), fmt.Errorf("unknown field object type %T", fType)
	}

	return sb.String(), nil
}

// ftype(b[0] | b[1] << 8 | b[2] << 16 ...)
func tpl_BytesToSimpleType(bufName string, fType *domain.SimpleFieldType, dst *fstringBuilder) {
	uintTypeName := fType.Name()
	if !fType.IsInteger() {
		uintTypeName = fmt.Sprintf("uint%d", fType.BitSize())
	}

	if !fType.IsInteger() {
		dst.WriteFString("%s(", fType.Name())
	}

	for i := 0; i < fType.Size(); i++ {
		dst.WriteFString("%s(%s[p", uintTypeName, bufName)
		if i != 0 {
			dst.WriteFString("+%d", i)
		}
		dst.WriteString("])")
		if i != 0 {
			dst.WriteFString("<<%d", i*8)
		}
		if i != fType.Size()-1 {
			dst.WriteString(" | ")
		}
		if i != 0 && i%4 == 0 {
			dst.WriteString("\n")
		}
	}
	if !fType.IsInteger() {
		dst.WriteString(")")
	}

	dst.WriteString("\n")
	dst.WriteFString("p += %d", fType.Size())
}
