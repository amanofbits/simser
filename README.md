# Simser
**Sim**ple **ser**ialization code generator for GO structs with unexported fields.  

Like `binary.Read`, but a bit different and without reflection.

## Features

### Core
- go module mode only
- simple sequential [de]serialization of simple structs
- supports non-exported fields, as well as exported
- basic (`int32`, `float64`, etc.), named (`type My uint64`, etc.) field types, and arrays of them (struct type fields,
array of arrays and slice of arrays are not supported right now)
- no reflection in generated code, it is simple and fast
- possibility to select type[s] to serialize via `-types` CLI flag
- customize output function names
- customize output file name

### Advanced
- slice field length can be set to any expression that returns `int`, by using tags. Currently [de]serialized instance can be referred to as "`o`", within expression.  
E.g. `simser:"len=o.PreviousIntegerField-5"`.  
Or `simser:"len=otherFunc()"`  
Remember that only fields that get read _before_ the slice field will have meaningful values (unless some tricks were used)

## Usage

#### Basic

`//go:generate go run github.com/am4n0w4r/simser -types=Header,body`
`//go:generate go run github.com/am4n0w4r/simser -types=all`

- `-types` (required): can be a comma-separated list of types you want to process, or reserved keyword `all` for processing all top-level `struct`s, found in file.  
`all` usage and requiredness of the argument can change in the future.
- `-output` (optional): set output file name.

#### Custom

`//go:generate go run github.com/am4n0w4r/simser -types=Header -output=file.name -read-fn-name=customReadFnName -write-fn-name=CustomWriteFnName`

- `-read-fn-name` (optional): custom name for deserializing function. Is set per-file.
- `-write-fn-name` (optional): custom name for deserializing function. Is set per-file.


## Project state

A bit messy, not very optimal, but simple and working. It was developed quickly from scratch, to serve a particular practical purpose, so the code itself is rather not perfect, but generated code should be good and do the job.  
It's the first time I worked with go's ast, so it was a lot of try-and-fail behind the scenes. Feel free to file issues.