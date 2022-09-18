namespace go sh.batch.schema

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
}