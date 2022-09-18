# Thrifty
---
Thrifty is used to transform a wire-format thrift message into a JSON representation
using an IDL definition. This library builds upon github.com/thrift-iterator/go by utilizing
the IDL definition in order to properly represent field names and enum values in the output
instead of IDs.


## Usage

First parse the IDL:

```go
idlFiles := map[string][]byte{
	"my.thrift": []byte(`
            namespace go sh.batch.schema
            
            struct Account {
              1: i32 id
              2: string first_name
              3: string last_name
              4: string email
            }
        `)
}


idl, err := thrifty.ParseIDLFiles(idlFiles)
if err != nil {
	log.Fatalf("unable to parse IDL files: %s", err.Error())
}
```

Then decode the binary message data. The struct name must be prefixed with the full namespace
 `"sh.batch.schema.Account"` as shown: 

```go
decodedMsg, err := thrify.DecodeWithParsedIDL(idlFiles, msgData, "sh.batch.schema.Account")
if err != nil {
	log.Fatalf("unable to decode thrift message: %s", err.Error())
}
```

The result of decodedMsg will be JSON in `[]byte` format:

```go
println(string(decodedMsg))
```

Output:
```json
{"first_name": "Gopher", "last_name": "Golang", "email": "gopher@golang.org", "id": 348590795}
```
