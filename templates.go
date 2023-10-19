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
	"go/types"
	"strings"
)

func tplReadField(f fieldInfo, bufName string) (t string, err error) {
	sb := strings.Builder{}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("/* %s */\n", f.name))
	sb.WriteString(tplReadNbytes(bufName, 1))
	sb.WriteByte('\n')

	switch f.typeToRead.Kind() {
	case types.Int8, types.Uint8,
		types.Int16, types.Uint16,
		types.Int32, types.Uint32,
		types.Int64, types.Uint64,
		types.Float32, types.Float64:
		s := tplConvertBytes(f.size, bufName, f.name, f.typeToAssign.String())
		sb.WriteString(s)
	default:
		return "", fmt.Errorf("no template found for %s", f.name)
	}
	sb.WriteByte('\n')
	return sb.String(), nil
}

func tplReadNbytes(bufName string, n uint) string {
	return fmt.Sprintf(
		`nRead, err := r.Read(%s)
n += int64(nRead)
if err != nil {
	return n, err
}
if nRead != %d {
	return n, err
}`, bufName, n)
}

func tplConvertBytes(size int, bufName, fieldName, typName string) string {
	utypName := typName
	isInteger := strings.HasPrefix(typName, "int") || strings.HasPrefix(typName, "uint")
	if !isInteger {
		utypName = fmt.Sprintf("uint%d", size*8)
	}

	sb := strings.Builder{}
	if size > 1 {
		sb.WriteString(fmt.Sprintf("_ = %s[%d]\n", bufName, size))
	}
	if isInteger {
		sb.WriteString(fmt.Sprintf("o.%s = ", fieldName))
	} else {
		sb.WriteString(fmt.Sprintf("o.%s = %s(", fieldName, typName))
	}

	for i := 0; i < size; i++ {
		sb.WriteString(fmt.Sprintf("%s(%s[%d])", utypName, bufName, i))
		if i != 0 {
			sb.WriteString(fmt.Sprintf("<<%d", i*8))
		}
		if i != size-1 {
			sb.WriteString(" | ")
		}
		if i != 0 && i%4 == 0 {
			sb.WriteString("\n")
		}
	}
	if !isInteger {
		sb.WriteString(")")
	}

	return sb.String()
}

func tplWriteField(f fieldInfo, bufName string) (t string, err error) {
	isInteger := strings.HasPrefix(f.typeToRead.Name(), "int") || strings.HasPrefix(f.typeToRead.Name(), "uint")
	intStr := fmt.Sprintf("uint%d", f.size*8)

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s = append(%s, ", bufName, bufName))
	for i := 0; i < f.size; i++ {
		sb.WriteString("byte(")
		if !isInteger {
			sb.WriteString(fmt.Sprintf("%s(", intStr))
		}
		sb.WriteString(fmt.Sprintf("s.%s", f.name))
		if !isInteger {
			sb.WriteString(")")
		}
		if i != 0 {
			sb.WriteString(fmt.Sprintf("<<%d", i*8))
		}
		sb.WriteString(")")
		if i != f.size-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")\n")
	return sb.String(), nil
}
