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
	"fmt"
	"strings"
)

// String builder which accepts formatted strings.
type fstringBuilder strings.Builder

func (b *fstringBuilder) WriteFString(frmt string, args ...any) {
	(*strings.Builder)(b).WriteString(fmt.Sprintf(frmt, args...))
}

func (b *fstringBuilder) WriteString(s string) (int, error) {
	return (*strings.Builder)(b).WriteString(s)
}

func (b *fstringBuilder) WriteByte(c byte) error {
	return (*strings.Builder)(b).WriteByte(c)
}

func (b *fstringBuilder) String() string { return (*strings.Builder)(b).String() }

func (b *fstringBuilder) Raw() *strings.Builder { return (*strings.Builder)(b) }
