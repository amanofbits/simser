# simser
**sim**ple struct binary **ser**ialization code generator for GO

This is a generator that should help with simple struct reading/writing to/from binary file.  
`binary.Read` does this well, but have a drawback - it doesn't work with non-exported fields.  
It also, as a runtime handler, has to do some checks to see how to do \[de]serialization.

## Project state

Messy, unoptimal but working. It was developed quickly from scratch, to serve a purpose, so the code itself is rather ugly, but the code it generates, should be at least ok.  
It's the first time I saw go's ast, so... Feel free to file issues.