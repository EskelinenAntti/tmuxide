package herdr

import (
	"os"
	"testing"
)

func TestInSession(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value string
		want  bool
	}{
		{name: "unset"},
		{name: "empty"},
		{name: "inactive", value: "0"},
		{name: "noncanonical", value: "true"},
		{name: "active", value: "1", want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HERDR_ENV", tt.value)
			if tt.name == "unset" {
				if err := os.Unsetenv("HERDR_ENV"); err != nil {
					t.Fatal(err)
				}
			}
			if got := InSession(); got != tt.want {
				t.Fatalf("InSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
