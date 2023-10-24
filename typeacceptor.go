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
	"errors"
	"fmt"
	"go/ast"

	"github.com/am4n0w4r/simser/internal/domain"
	"golang.org/x/exp/maps"
)

type typeAcceptor map[string]uint8

func newTypeAcceptor(rawTypeList []string) (*typeAcceptor, error) {
	if len(rawTypeList) == 0 {
		return nil, errors.New("no input types")
	}
	st := map[string]uint8{}
	for _, v := range rawTypeList {
		if v == "" {
			return nil, errors.Join(domain.ErrUnsupportedType, errors.New("empty type name in type list"))
		}
		st[v] = 0
	}
	return (*typeAcceptor)(&st), nil
}

func (st *typeAcceptor) Accepts(ts *ast.TypeSpec) bool {
	if _, all := (*st)["all"]; len(*st) == 1 && all {
		return true
	}
	_, ok := (*st)[ts.Name.Name]
	if ok {
		delete(*st, ts.Name.Name)
	}
	return ok
}

func (ts typeAcceptor) IsDrained() bool {
	return len(ts) == 0
}

func (ta typeAcceptor) AcceptsAll() bool {
	_, hasAll := ta["all"]
	return len(ta) == 1 && hasAll
}

func (ts typeAcceptor) String() string { return fmt.Sprintf("%s", maps.Keys(ts)) }
