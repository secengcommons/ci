package tools

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/secengcommons/verify/goverify"
)

func TestWorkflowsUseApprovedImmutableReferences(t *testing.T) {
	directory := filepath.Join("..", ".github", "workflows")
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if closeErr := root.Close(); closeErr != nil {
			t.Error(closeErr)
		}
	})
	found := false
	for _, entry := range entries {
		if entry.IsDir() || !workflowExtension(entry.Name()) {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			t.Fatalf("workflow %s is a symbolic link", entry.Name())
		}
		found = true
		checkWorkflow(t, root, entry.Name())
	}
	if !found {
		t.Fatal("workflow inventory is empty")
	}
}

func TestWorkflowReferenceRegressions(t *testing.T) {
	for name, reference := range map[string]string{
		"branch":          "actions/checkout@main # v7.0.1",
		"tag":             "actions/checkout@v7 # v7.0.1",
		"short":           "actions/checkout@3d3c42e # v7.0.1",
		"unknown":         "example/action@" + strings.Repeat("a", 40) + " # v1.0.0",
		"missing version": "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1",
		"wrong version":   "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.0",
		"local action":    "./.github/actions/check",
	} {
		t.Run(name, func(t *testing.T) {
			source := []byte("name: Invalid\non: push\njobs:\n  check:\n    runs-on: ubuntu-24.04\n    steps:\n      - uses: " + reference + "\n")
			references, err := goverify.InspectWorkflowReferences(source)
			if err == nil && approvedReferences(references) {
				t.Fatal("invalid workflow reference was accepted")
			}
		})
	}
}

func TestWorkflowsUseRepositoryVerifyTool(t *testing.T) {
	expected := map[string][]string{
		"alpine.yml": {"go tool secverify test"},
		"fuzz.yml":   {"go tool secverify campaign"},
		"linux.yml":  {"go tool secverify test"},
		"macos.yml":  {"go tool secverify test"},
		"verify.yml": {"go tool secverify all"},
		"windows.yml": {
			"go tool secverify test",
			"go tool secverify fuzz-inventory",
		},
	}
	for path, commands := range expected {
		source, err := os.ReadFile(filepath.Join("..", ".github", "workflows", path))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(source, []byte("go install github.com/secengcommons/verify")) ||
			bytes.Contains(source, []byte("run: secverify")) || bytes.Contains(source, []byte("/bin/secverify")) {
			t.Fatalf("%s uses an external Verify executable", path)
		}
		for _, command := range commands {
			if bytes.Count(source, []byte(command)) != 1 {
				t.Fatalf("%s does not own exactly one %q command", path, command)
			}
		}
		if !bytes.Contains(source, []byte("go run -a -buildvcs=true ./cmd/secverify")) {
			t.Fatalf("%s does not retain Verify source bootstrap", path)
		}
	}
}

func checkWorkflow(t *testing.T, root *os.Root, path string) {
	t.Helper()
	file, err := root.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	information, statErr := file.Stat()
	if statErr != nil {
		t.Fatal(errors.Join(statErr, file.Close()))
	}
	if information == nil {
		if closeErr := file.Close(); closeErr != nil {
			t.Error(closeErr)
		}
		t.Fatal("workflow information is absent")
	}
	if information.Size() < 0 || information.Size() > goverify.MaxWorkflowBytes {
		if closeErr := file.Close(); closeErr != nil {
			t.Error(closeErr)
		}
		t.Fatalf("%s exceeds its byte bound", path)
	}
	source, readErr := io.ReadAll(io.LimitReader(file, goverify.MaxWorkflowBytes+1))
	if err = errors.Join(readErr, file.Close()); err != nil {
		t.Fatal(err)
	}
	references, err := goverify.InspectWorkflowReferences(source)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	if !approvedReferences(references) {
		t.Fatalf("%s contains an unapproved action", path)
	}
}

func approvedReferences(references []goverify.WorkflowReference) bool {
	approved := map[string]string{
		"actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1":                 "v7.0.1",
		"actions/dependency-review-action@a1d282b36b6f3519aa1f3fc636f609c47dddb294": "v5.0.0",
		"actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e":                 "v7.0.0",
		"cross-platform-actions/action@faa0c6197e94aacf1c5956460152c8380d3560a5":    "v1.5.0",
		"github/codeql-action/analyze@cdf488f595d80d6e07e03d4674febd5ab45fa938":     "v4.37.9",
		"github/codeql-action/init@cdf488f595d80d6e07e03d4674febd5ab45fa938":        "v4.37.9",
		"vmactions/alpine-vm@697847385681cdc413604693010eded91e8d9650":              "v1.0.1",
		"vmactions/solaris-vm@96d8d976f9e67d82ec6c7e8ce9c1060731f9e21c":             "v1.3.9",
	}
	for _, reference := range references {
		if reference.Kind != goverify.WorkflowAction || strings.HasPrefix(reference.Value, "docker://") {
			continue
		}
		if strings.HasPrefix(reference.Value, "./.github/workflows/") {
			continue
		}
		version, found := approved[reference.Value]
		if !found || version != reference.Version {
			return false
		}
	}
	return true
}

func workflowExtension(path string) bool {
	extension := filepath.Ext(path)
	return extension == ".yml" || extension == ".yaml"
}
