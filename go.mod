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

module github.com/amanofbits/simser

// Before Go 1.21, the directive was advisory only; now it is a mandatory requirement
// Source: https://go.dev/doc/modules/gomod-ref#go
go 1.18

require (
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	golang.org/x/mod v0.13.0
	golang.org/x/tools v0.14.0
)

require golang.org/x/sys v0.13.0 // indirect
