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

package generator

import (
	"github.com/am4n0w4r/simser/internal/domain"
)

func GenStructCode(s domain.InputStruct, out *Output, readFnName, writeFnName string) error {
	if s.FieldCount() == 0 {
		return nil
	}

	totalBufLen := 0
	for i := 0; i < s.FieldCount(); i++ {
		totalBufLen += s.Field(i).Type().Size()
	}

	// Generate func (o 'typename')LoadFrom(io.Reader) (*typename, error)
	{
		out.AppendImport("io")

		out.AppendF("func (o *%s) %s(r io.Reader) (n int, err error) {\n", s.Name(), readFnName)
		out.AppendF("b := make([]byte, %d)\n", totalBufLen)
		out.AppendF("p := 0\n")
		out.LF()
		out.Append(tpl_ReadNbytesIntoBuf("b", uint(totalBufLen))).LF()
		out.LF()

		for i := 0; i < s.FieldCount(); i++ {
			field := s.Field(i)
			out.AppendF("// %s", field.Name()).LF()
			s, err := tpl_ReadField(field, "b")
			if err != nil {
				return err
			}
			out.AppendF("%s\n", s)
		}

		out.LF()
		out.Append("return n, err")
		out.Append("}\n")
	}

	out.LF()

	// Gen func (o 'typename') SaveTo(io.Writer) (n int, err error)
	{
		out.AppendImport("io")

		out.AppendF("func (s *%s) %s(w io.Writer) (n int, err error) {\n", s.Name(), writeFnName)
		out.AppendF("b := make([]byte, 0, %d)\n", totalBufLen)
		out.LF()

		for i := 0; i < s.FieldCount(); i++ {
			field := s.Field(i)
			out.AppendF("// %s", field.Name()).LF()
			s, err := tpl_WriteField(field, "b")
			if err != nil {
				return err
			}
			out.AppendF("%s\n", s)
		}

		out.Append("return w.Write(b)")
		out.Append("}\n")
	}

	return nil
}
