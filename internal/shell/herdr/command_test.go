package herdr

import (
	"os/exec"
	"strings"
	"testing"
)

func TestShellCommand(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want string
	}{
		{"no arguments", nil, ""},
		{"editor and flags", []string{"nvim", "-R", "/project/file"}, "'nvim' '-R' '/project/file'"},
		{"empty argument", []string{""}, "''"},
		{"spaces", []string{"my editor", "a b"}, "'my editor' 'a b'"},
		{"single quotes", []string{"a'b'c"}, "'a'\"'\"'b'\"'\"'c'"},
		{"shell syntax", []string{"$HOME; $(exit 1) `exit 1` * >file & |", "a\nb\t\\\""}, "'$HOME; $(exit 1) `exit 1` * >file & |' 'a\nb\t\\\"'"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellCommand(tt.args); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShellCommandPreservesArguments(t *testing.T) {
	args := []string{
		"", "plain", "a b", "'", "a'b'c", "\"double quotes\"", "back\\slash",
		"line\nbreak", "tab\there", "$HOME", "$(exit 1)", "`exit 1`",
		"; exit 1", "* ? [abc]", ">file & | ( )", "日本語",
	}
	// printf's NUL separators make empty arguments and embedded newlines observable.
	command := append([]string{"printf", "%s\\0"}, args...)
	output, err := exec.Command("sh", "-c", shellCommand(command)).Output()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(output), strings.Join(args, "\x00")+"\x00"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
