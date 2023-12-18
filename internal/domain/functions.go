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

package domain

import (
	"fmt"
	"strconv"
	"strings"
)

/******************************************************************/
/* This file should be kept as small as possible.                 */
/* Ony functions that are used in domain methods and ALSO outside */
/******************************************************************/

// Adds parentheses to Go's int expression if it is
// a) not a simple int value,
// b) not already cast to int.
func ParenthesizeIntExpr(expr string) string {
	if _, err := strconv.Atoi(expr); err == nil {
		return expr
	}
	if !strings.HasPrefix(expr, "int") || !strings.HasSuffix(expr, ")") {
		return fmt.Sprintf("(%s)", expr)
	}
	pos := 3
	for expr[pos] == ' ' {
		pos++
	}
	if expr[pos] != '(' {
		return fmt.Sprintf("(%s)", expr)
	}
	return expr
}

// Returns true if total argument's size can be interpreted as known at declaration time.
// E.g. (primitives, arrays). NOT slices.
func IsFixedSize[
	T StructField | SimpleFieldType | ArrayFieldType | SliceFieldType | int](a T) bool {

	switch arg := any(a).(type) {
	case StructField:
		return arg.Type().Size() >= 0
	case SimpleFieldType:
		return arg.Size() >= 0
	case ArrayFieldType:
		return arg.Size() >= 0
	case SliceFieldType:
		return arg.Size() >= 0
	case int:
		return arg >= 0
	default:
		panic(fmt.Errorf("unknown type %T", a))
	}
}
