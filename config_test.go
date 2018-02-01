package main

import (
	"os"
	"testing"
)

func TestGetenv(t *testing.T) {
	type checkFunc func(string) error

	tests := [...]struct {
		name   string
		key    string
		env    string
		def    string
		expect string
	}{
		{
			name:   "fetches from env",
			key:    "MY_TEST_KEY",
			env:    "theenvvalue",
			def:    "",
			expect: "theenvvalue",
		},
		{
			name:   "fetches from env even if default is set",
			key:    "MY_TEST_KEY",
			env:    "theenvvalue",
			def:    "thedefaultvalue",
			expect: "theenvvalue",
		},
		{
			name:   "uses defaults if env is empty",
			key:    "MY_TEST_KEY",
			env:    "",
			def:    "thedefaultvalue",
			expect: "thedefaultvalue",
		},
		{
			name:   "fetches weird values from env",
			key:    "the TEST key",
			env:    "this is the value. \n\n //\\\\ \nEOF\n Whynot",
			expect: "this is the value. \n\n //\\\\ \nEOF\n Whynot",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(tc.key, tc.env)
			have := getenv(tc.key, tc.def)

			if have != tc.expect {
				t.Errorf("expected value %q, found %q", tc.expect, have)
			}

			if err := os.Unsetenv(tc.key); err != nil {
				t.Fatalf("Unable to unset the key %q: %v", tc.key, err)
			}
		})
	}
}
