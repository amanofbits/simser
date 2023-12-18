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

func GenStructCode(s domain.InputStruct, out *Output, readFnName, writeFnName string) error {
	if s.FieldCount() == 0 {
		return nil
	}

	sizeGroups := getFieldSizeGroups(s)

	// Generate func (o 'typename')LoadFrom(io.Reader) (*'typename', error)
	{
		out.AppendImport("io")

		out.AppendF("func (o *%s) %s(r io.Reader) (n int, err error) {\n", s.Name(), readFnName)
		out.AppendF("p, nRead := 0, 0\n")

		out.AppendF("toRead := ")
		if domain.IsFixedSize(sizeGroups[0]) {
			out.AppendF("%d\n", sizeGroups[0])
		} else {
			out.AppendF("%s\n", s.Field(0).Type().SizeExpr())
		}

		for i := 0; i < s.FieldCount(); i++ {
			if s.Field(i).Type().IsSequence() {
				out.AppendF("sLen, sElSize := 0, 0\n")
				break
			}
		}

		out.LF()
		out.AppendF("b := make([]byte, toRead)\n")
		out.Append(tpl_ReadBytesIntoBuf("b")).LF()
		out.LF()

		for i := 0; i < s.FieldCount(); i++ {
			field := s.Field(i)
			out.AppendF("\n// %s\n", field.Name())
			if i != 0 {
				if size, ok := sizeGroups[i]; ok {
					if !domain.IsFixedSize(size) {
						if field.Type().IsSequence() {
							seqType := field.Type().(domain.SequenceFieldType)
							out.AppendF("sLen, sElSize = %s, %s\n", seqType.LenExpr(), seqType.ElType().SizeExpr())
						}
						out.Append("p, toRead = 0, sLen * sElSize\n")
					} else {
						out.AppendF("p, toRead = 0, %d\n", size)
					}
					out.AppendF("if toRead > cap(b) {\n")
					out.AppendF("b = make([]byte, toRead)\n")
					out.Append("}\n")
					out.Append(tpl_ReadBytesIntoBuf("b")).LF()
				}
			}
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

		out.AppendF("func (o *%s) %s(w io.Writer) (n int, err error) {\n", s.Name(), writeFnName)

		out.AppendF("b := make([]byte, 0, ")
		sb := fstringBuilder{}
		constSize := 0
		for i, size := range sizeGroups {
			if !domain.IsFixedSize(size) {
				expr := s.Field(i).Type().SizeExpr()
				sb.WriteFString("+ %s", domain.ParenthesizeIntExpr(expr))
				continue
			}
			constSize += size
		}
		out.AppendF("%d %s)\n", constSize, sb.String())
		out.LF()

		for i := 0; i < s.FieldCount(); i++ {
			field := s.Field(i)
			out.AppendF("\n// %s", field.Name()).LF()
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
