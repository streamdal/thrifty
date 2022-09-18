package thrifty

import "testing"

func TestGetMessageAndNS(t *testing.T) {
	cases := []struct {
		In         string
		Namespace  string
		StructName string
	}{
		{
			"included.Account",
			"included",
			"Account",
		},
	}

	for _, testCase := range cases {
		name, ns, err := getMessageAndNS(testCase.In)
		if err != nil {
			t.Error(err)
			continue
		}

		if ns != testCase.Namespace {
			t.Errorf("case = %s, wanted = %s, got = %s", testCase.In, testCase.Namespace, ns)
		}

		if name != testCase.StructName {
			t.Errorf("case = %s, wanted = %s, got = %s", testCase.In, testCase.StructName, name)
		}
	}
}

func TestGetMessageAndNS_Error(t *testing.T) {
	_, _, err := getMessageAndNS("Account")
	if err == nil {
		t.Fail()
	}
}

// TODO: add tests
