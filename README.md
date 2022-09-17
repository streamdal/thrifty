# Thrifty
---
Thrifty is used to transform a wire-format thrift message into a JSON representation
using an IDL definition. This library builds upon github.com/thrift-iterator/go by utilizing
the IDL definition in order to properly represent field names and enum values in the output
instead of IDs.
