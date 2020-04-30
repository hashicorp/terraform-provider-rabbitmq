package rabbitmq

import "testing"

func TestParseId(t *testing.T) {
	var badInputs = []string{
		"",
		"foo/test",
		"footest",
		"foo@bar@test",
	}

	for _, input := range badInputs {
		_, _, err := parseId(input)
		if err == nil {
			t.Errorf("parseId failed for: %s.", input)
		}
	}

	var goodInputs = []struct {
		input string
		name  string
		vhost string
	}{
		{"foo@test", "foo", "test"},
		{"foo@/", "foo", "/"},
		{"foo/bar/baz@/", "foo/bar/baz", "/"},
	}

	for _, test := range goodInputs {
		name, vhost, err := parseId(test.input)
		if err != nil || name != test.name || vhost != test.vhost {
			t.Errorf("parseId failed for: %s.", test.input)
		}
	}
}
