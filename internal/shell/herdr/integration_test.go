package herdr

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/shell/runner"
)

// Opt in with HERDR_TEST_BIN=/path/to/herdr go test ./internal/shell/herdr.
// Everything, including the plugin registry, lives in an isolated temp directory.
func TestNativeHerdr(t *testing.T) {
	binary := os.Getenv("HERDR_TEST_BIN")
	if binary == "" {
		t.Skip("set HERDR_TEST_BIN to run against Herdr 0.9.5+")
	}
	var err error
	binary, err = exec.LookPath(binary)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	t.Setenv("HERDR_CONFIG_PATH", filepath.Join(root, "config.toml"))
	t.Setenv("HERDR_SOCKET_PATH", filepath.Join(root, "herdr.sock"))
	t.Setenv("HERDR_SESSION", "")
	t.Setenv("HERDR_ENV", "1")
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(root, "state"))
	t.Setenv("SHELL", "/bin/sh")
	log, err := os.Create(filepath.Join(root, "server.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	server := exec.Command(binary, "server")
	server.Stdout, server.Stderr = log, log
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- server.Wait() }()
	t.Cleanup(func() {
		stop := exec.Command(binary, "server", "stop")
		_ = stop.Run()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = server.Process.Kill()
			<-done
		}
	})
	host := Cmd{Runner: runner.CmdRunner{}, Binary: binary}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err = host.call("workspace", "list"); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("server did not become ready:", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	proj, err := project.ForDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := host.Open(proj, nil); err != nil {
			t.Fatal(err)
		}
	}
	list, err := host.call("workspace", "list")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	var id string
	for _, workspace := range list.Workspaces {
		if workspace.Label == proj.Name {
			count++
			id = workspace.ID
		}
	}
	if count != 1 {
		t.Fatalf("expected one reused workspace, got %d", count)
	}
	file := filepath.Join(root, "a'b; $(false).txt")
	if err := os.WriteFile(file, []byte("TMUXIDE_NATIVE_SMOKE_OK\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := host.Open(proj, []string{"cat", file}); err != nil {
		t.Fatal(err)
	}
	// The created editor tab must become the workspace's active tab.
	data, err := host.output("workspace", "get", id)
	if err != nil {
		t.Fatal(err)
	}
	var info struct {
		Result struct {
			Workspace struct {
				TabCount    int    `json:"tab_count"`
				ActiveTabID string `json:"active_tab_id"`
			} `json:"workspace"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &info); err != nil {
		t.Fatal(err)
	}
	if info.Result.Workspace.TabCount != 2 {
		t.Fatalf("expected a new editor tab: %s", data)
	}
	data, err = host.output("pane", "list", "--workspace", id)
	if err != nil {
		t.Fatal(err)
	}
	var panes struct {
		Result struct {
			Panes []struct {
				ID    string `json:"pane_id"`
				TabID string `json:"tab_id"`
			} `json:"panes"`
		} `json:"result"`
	}
	if err := json.Unmarshal(data, &panes); err != nil {
		t.Fatal(err)
	}
	var paneID string
	for _, pane := range panes.Result.Panes {
		if pane.TabID == info.Result.Workspace.ActiveTabID {
			paneID = pane.ID
		}
	}
	if paneID == "" {
		t.Fatal("editor tab has no pane")
	}
	// Check real PTY output, not just successful command queueing.
	deadline = time.Now().Add(5 * time.Second)
	for {
		data, err = host.output("pane", "read", paneID)
		if err == nil && strings.Contains(string(data), "TMUXIDE_NATIVE_SMOKE_OK") {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("editor command did not run: %s (%v)", data, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	repo, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := host.output("plugin", "link", repo); err != nil {
		t.Fatal(err)
	}
}
