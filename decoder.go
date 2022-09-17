// Package thrifty is used to transform a wire-format thrift message into a JSON representation
// using an IDL definition. This library builds upon github.com/thrift-iterator/go by utilizing
// the IDL definition in order to properly represent field names and enum values in the output
// instead of IDs.
package thrifty

import (
	"fmt"
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	thrifter "github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/general"
	"github.com/thrift-iterator/go/protocol"
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
)

// ParsedIDL is a convenience struct that holds AST representations of Thrift structs
// and mappings of const and enum types.
type ParsedIDL struct {
	Structs map[string]*ast.Struct
	Enums   map[string]map[int32]string
}

// DecodeWithParsedIDL decodes a thrift message into JSON format using an IDL abstract syntax tree.
// It is recommended to use this method instead of DecodeWithRawIDL if you have multiple messages to
// decode using the same IDL. Before calling this method, you must parse the IDL definition using ParseIDL
func DecodeWithParsedIDL(idl *ParsedIDL, thriftMsg []byte, structName string) ([]byte, error) {
	decoded, err := decodeWireFormat(thriftMsg)
	if err != nil {
		return nil, err
	}

	result, err := structToMap(idl, structName, decoded)
	if err != nil {
		return nil, err
	}

	// jsoniter is needed to marshal map[interface{}]interface{} types
	js, err := jsoniter.Marshal(result)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal decoded thrift message to JSON")
	}

	return js, nil
}

// DecodeWithRawIDL decodes a thrift message with the provided IDL definition and root message name
// It is recommended to use DecodeWithParsedIDL() instead to avoid the overhead of having to parse the IDL
// into an AST on every call. This method is here for convenience purposes
func DecodeWithRawIDL(idlDefinition []byte, thriftMsg []byte, structName string) ([]byte, error) {
	idl, err := ParseIDL(idlDefinition)
	if err != nil {
		return nil, err
	}

	return DecodeWithParsedIDL(idl, thriftMsg, structName)
}

// ParseIDL takes an IDL definition and returns a ParsedIDL struct containing AST representations
// of all thrift structs, and a mapping of enum int->string values. All other Thrift IDL
func ParseIDL(data []byte) (*ParsedIDL, error) {
	parsedIDL := &ParsedIDL{
		Structs: make(map[string]*ast.Struct),
		Enums:   make(map[string]map[int32]string),
	}

	parsed, err := idl.Parse(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse IDL")
	}

	for _, def := range parsed.Definitions {
		constant, ok := def.(*ast.Constant)
		if ok {
			// TODO: handle constants
			_ = constant
		}

		enums, ok := def.(*ast.Enum)
		if ok {
			parsedIDL.Enums[enums.Name] = make(map[int32]string)
			for _, enum := range enums.Items {
				if enum.Value == nil {
					// TODO: enums without values pass IDL parser, I'm guessing they default to iota style?
					continue
				}

				parsedIDL.Enums[enums.Name][int32(*enum.Value)] = enum.Name
			}
		}

		msg, ok := def.(*ast.Struct)
		if !ok {
			// Ignore non structs
			continue
		}

		parsedIDL.Structs[msg.Name] = msg
	}

	return parsedIDL, nil
}

func decodeWireFormat(message []byte) (*general.Struct, error) {
	obj := &general.Struct{}

	err := thrifter.Unmarshal(message, obj)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read thrift message")
	}

	return obj, nil
}

func structToMap(idl *ParsedIDL, rootMsgName string, decoded *general.Struct) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})

	rootMsg, ok := idl.Structs[rootMsgName]
	if !ok {
		log.Printf("message '%s' not found in IDL", rootMsgName)
		return jsonMap, nil
	}

	for _, field := range rootMsg.Fields {
		// Non-base type
		if _, ok := field.Type.(ast.TypeReference); ok {
			// Check if field is a constant or enum
			enums, ok := idl.Enums[field.Type.String()]
			if ok {
				enumID, ok := decoded.Get(protocol.FieldId(field.ID)).(int32)
				if !ok {
					return nil, fmt.Errorf("could not type assert ID for field '%s' to int32", field.Name)
				}

				jsonMap[field.Name] = enums[enumID]

				continue
			}

			// Field IDs can be repeated between structs. Recurse down the decoded data
			subType, ok := decoded.Get(protocol.FieldId(field.ID)).(general.Struct)
			if !ok {
				return nil, fmt.Errorf("could not type assert field '%s' to general.Struct", field.Name)
			}

			v, err := structToMap(idl, field.Type.String(), &subType)
			if err != nil {
				return nil, err
			}

			jsonMap[field.Name] = v

			continue

		}

		// Scalar type
		fieldVal := decoded.Get(protocol.FieldId(field.ID))
		jsonMap[field.Name] = fieldVal
	}

	return jsonMap, nil
}
