//go:build mage
// +build mage

package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
)

func init() {
}

func Deps() {
	deps.GetAllDeps()
}

// Build binaries.
func Build() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg := path.Join(cwd, ".goreleaser.yml")
	return common.BuildReleaser("--config", cfg, "--clean")
}

// Release releases the project.
func Release() error {
	return common.Release()
}

func GetWorkspacePaths() []string {
	return []string{"."}
}

// Test - based on ci.yaml implementation:
// go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} .ext/gobin/gotestsum-v1.10.0/gotestsum --format short-verbose -- -count=1 -v {}/...
func Test() error {
	for _, p := range GetWorkspacePaths() {
		if err := deps.GoDep("gotestsum")([]string{"--format", "short-verbose", "--", "-count=1", "-v", "./" + filepath.Join(p, "...")}...); err != nil {
			return err
		}
	}
	return nil
}

// Lint - based on ci.yaml implementation:
// go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} .ext/gobin/golangci-lint-v1.52.2/golangci-lint run {}/... -c .golangci.yaml
func Lint() error {
	for _, p := range GetWorkspacePaths() {
		if err := deps.GoDep("golangci-lint")([]string{"run", "./" + filepath.Join(p, "..."), "-c", ".golangci.yaml"}...); err != nil {
			return err
		}
	}
	return nil
}
