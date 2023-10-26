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
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

type structTag struct {
	values map[string]string
}

func newStructTag() *structTag {
	return &structTag{}
}

// Deprecated: Use specific methods to get specific values.
// This note will be here as a reminder.
// func (p structTag) get(key string) (val string, ok bool) {
// 	val, ok = p.values[key]
// 	return val, ok
// }

func (p *structTag) parse(t *ast.BasicLit) error {
	p.values = map[string]string{}
	if t == nil {
		return nil
	}
	if t.Kind != token.STRING {
		return fmt.Errorf("illegal tag kind %s", t.Kind)
	}
	unqTag, err := strconv.Unquote(t.Value)
	if err != nil {
		return err
	}
	tag := reflect.StructTag(unqTag) // reflect has "canonical" parser for tags

	simser, ok := tag.Lookup("simser")
	if !ok {
		return nil
	}

	vals := strings.Split(simser, ",")
	for _, v := range vals {
		idx := strings.Index(v, "=")
		if idx < 0 {
			p.values[v] = ""
		}
		p.values[v[:idx]] = v[idx+1:]
	}
	return nil
}

func (p structTag) getLenExpr() (expr string, ok bool, err error) {
	key := "len"
	expr, ok = p.values[key]
	if !ok {
		return expr, ok, nil
	}
	return expr, ok, p._validateExpr(key, expr)
}

var ErrCommentInTagExpr = errors.New("comments are not allowed within tag expressions")

func (p structTag) _validateExpr(key, expr string) error {
	// Comments can confuse further parsing, they are not allowed within tags.
	// Not that anyone would put them there, but...
	commentStart, commentEnd := false, false // "//", "/*" ; "*/"
	for i, r := range expr {
		if commentStart {
			if r == '/' || r == '*' {
				return errors.Join(ErrCommentInTagExpr, fmt.Errorf("'%s', at char %d", key, i-1))
			}
			commentStart = false
		} else if commentEnd {
			if r == '/' {
				return errors.Join(ErrCommentInTagExpr, fmt.Errorf("'%s', at char %d", key, i-1))
			}
			commentEnd = false
		}
		if r == '/' {
			commentStart = true
		} else if r == '*' {
			commentEnd = true
		}
	}
	return nil
}
