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
	"github.com/amanofbits/simser/internal/domain"
)

// Returns a description of size groups as '[index]size', where:
// a) sizes of all consecutive fixed-size fields are added and put with index of the first field in group.
// b) variable-sized fields are never grouped (all their indices present in output) and their size is always -1
func getFieldSizeGroups(s domain.InputStruct) (g map[int]int) {
	g = map[int]int{}
	if s.FieldCount() < 1 {
		return g
	}

	sizes := make([]int, s.FieldCount())
	for i := 0; i < s.FieldCount(); i++ {
		sizes[i] = -1
		if domain.IsFixedSize(s.Field(i)) {
			sizes[i] = s.Field(i).Type().Size()
		}
	}

	sum, startIdx := 0, -1
	for i := 0; i < s.FieldCount(); i++ {
		fsize := s.Field(i).Type().Size()
		if !domain.IsFixedSize(s.Field(i)) {
			if i != startIdx {
				g[startIdx] = sum
				sum = 0
			}
			startIdx = i + 1
			g[i] = fsize
			continue
		}
		if startIdx < 0 {
			startIdx = i
		}
		sum += fsize
	}
	if sum > 0 {
		g[startIdx] = sum
	}

	return g
}
