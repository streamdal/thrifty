package thrifty

import (
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
		It("Parses a single file", func() {
			data := []byte(`namespace go sh.batch.schema

include "somefile.thrift"

struct SubMessage {
    1: string value
}

enum ClientType {
  UNSET = 0,
  VIP = 1
}

union Thing {
  1: string thing_string
  2: i32 thing_int
}

const i32 INT_CONST = 1234;

typedef double USD

struct Customer {
  1: i32 key
  2: string value
  3: SubMessage subm
  4: map<i32, i32> newmap
  5: list<string> newlist
  6: ClientType client_type = ClientType.VIP
  7: Thing unionthing
  8: USD monthly_price
  9: i32 testconst = INT_CONST
}`)

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
})
