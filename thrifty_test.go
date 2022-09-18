package thrifty

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Thrifty", func() {
	Context("parseStructName", func() {
		It("returns correct namespace and struct name", func() {
			name, ns, err := parseStructName("incfile.Account")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("Account"))
			Expect(ns).To(Equal("incfile"))
		})

		It("returns an error when no namespace is present", func() {
			_, _, err := parseStructName("Account")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("'Account' must contain a namespace"))
		})
	})

	Context("ParseIDL", func() {
		It("Parses a definition", func() {
			data, err := os.ReadFile("./test-assets/complex.thrift")

			idl, err := ParseIDL(data)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(idl.Typedefs)).To(Equal(1))
			Expect(idl.Typedefs).To(HaveKey("USD"))
			Expect(len(idl.Structs)).To(Equal(3))
			Expect(len(idl.Enums)).To(Equal(1))
			Expect(idl.Enums["ClientType"][0]).To(Equal("UNSET"))
			Expect(idl.Namespace).To(Equal("sh.batch.schema"))
		})
	})

	Context("DecodeWithParsedIDL", func() {
		It("decodes a simple message into JSON", func() {
			msgData, err := os.ReadFile("./test-assets/simple.bin")
			Expect(err).ToNot(HaveOccurred())

			idlData, err := os.ReadFile("./test-assets/simple.thrift")
			Expect(err).ToNot(HaveOccurred())

			idl, err := ParseIDLFiles(map[string][]byte{"simple.thrift": idlData})
			Expect(err).ToNot(HaveOccurred())

			decoded, err := DecodeWithParsedIDL(idl, msgData, "sh.batch.schema.Account")
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(MatchJSON(`{"id": 348590795, "first_name": "Gopher", "last_name": "Golang", "email": "gopher@golang.org"}`))
		})

		It("decodes a nested message into JSON", func() {
			expected, err := os.ReadFile("./test-assets/nested.json")
			Expect(err).ToNot(HaveOccurred())

			msgData, err := os.ReadFile("./test-assets/nested.bin")
			Expect(err).ToNot(HaveOccurred())

			idlData, err := os.ReadFile("./test-assets/nested.thrift")
			Expect(err).ToNot(HaveOccurred())

			idl, err := ParseIDLFiles(map[string][]byte{"nested.thrift": idlData})
			Expect(err).ToNot(HaveOccurred())

			decoded, err := DecodeWithParsedIDL(idl, msgData, "sh.batch.schema.Account")
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(MatchJSON(expected))
		})
	})
})

func BenchmarkParseIDL_simple(b *testing.B) {
	idlData, err := os.ReadFile("./test-assets/simple.thrift")
	Expect(err).ToNot(HaveOccurred())

	idlFiles := map[string][]byte{"simple.thrift": idlData}

	for n := 0; n < b.N; n++ {
		ParseIDLFiles(idlFiles)
	}
}

func BenchmarkDecodeWithParsedIDL_simple(b *testing.B) {
	msgData, err := os.ReadFile("./test-assets/simple.bin")
	Expect(err).ToNot(HaveOccurred())

	idlData, err := os.ReadFile("./test-assets/simple.thrift")
	Expect(err).ToNot(HaveOccurred())

	idl, err := ParseIDLFiles(map[string][]byte{"simple.thrift": idlData})
	Expect(err).ToNot(HaveOccurred())

	for n := 0; n < b.N; n++ {
		DecodeWithParsedIDL(idl, msgData, "sh.batch.schema.Account")
	}
}

func BenchmarkDecodeWithRawIDL_simple(b *testing.B) {
	msgData, err := os.ReadFile("./test-assets/simple.bin")
	Expect(err).ToNot(HaveOccurred())

	idlData, err := os.ReadFile("./test-assets/simple.thrift")
	Expect(err).ToNot(HaveOccurred())

	idlFiles := map[string][]byte{"simple.thrift": idlData}

	for n := 0; n < b.N; n++ {
		DecodeWithRawIDL(idlFiles, msgData, "sh.batch.schema.Account")
	}
}

func BenchmarkParseIDL_nested(b *testing.B) {
	idlData, err := os.ReadFile("./test-assets/nested.thrift")
	Expect(err).ToNot(HaveOccurred())

	idlFiles := map[string][]byte{"nested.thrift": idlData}

	for n := 0; n < b.N; n++ {
		ParseIDLFiles(idlFiles)
	}
}

func BenchmarkDecodeWithParsedIDL_nested(b *testing.B) {
	msgData, err := os.ReadFile("./test-assets/nested.bin")
	Expect(err).ToNot(HaveOccurred())

	idlData, err := os.ReadFile("./test-assets/nested.thrift")
	Expect(err).ToNot(HaveOccurred())

	idl, err := ParseIDLFiles(map[string][]byte{"nested.thrift": idlData})
	Expect(err).ToNot(HaveOccurred())

	for n := 0; n < b.N; n++ {
		DecodeWithParsedIDL(idl, msgData, "sh.batch.schema.Account")
	}
}

func BenchmarkDecodeWithRawIDL_nested(b *testing.B) {
	msgData, err := os.ReadFile("./test-assets/nested.bin")
	Expect(err).ToNot(HaveOccurred())

	idlData, err := os.ReadFile("./test-assets/simple.thrift")
	Expect(err).ToNot(HaveOccurred())

	idlFiles := map[string][]byte{"simple.thrift": idlData}

	for n := 0; n < b.N; n++ {
		DecodeWithRawIDL(idlFiles, msgData, "sh.batch.schema.Account")
	}
}
