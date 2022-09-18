namespace go sh.batch.schema

struct Account {
  1: i32 id
  2: string first_name
  3: string last_name
  4: string email
  5: Billing billing
  6: Address address
  7: Deep1 deep_nested
}

struct Billing {
  1: string card_number
  2: i32 exp_month
  3: i32 exp_year
}

struct Address {
  1: string street
  2: string city
  3: string state_province
  4: string country
  5: string postal_code
}

struct Deep1 {
  1: Deep2 deep2
}

struct Deep2 {
  1: Deep3 deep3
}

struct Deep3 {
  1: string nested_value
}