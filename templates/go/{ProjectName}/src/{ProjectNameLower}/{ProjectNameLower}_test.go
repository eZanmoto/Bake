// Copyright {Year} {Owner}. All rights reserved.

package main

{?bin}
import (
	"fmt"
	"os"
	"path"
	"testing"
)

var (
	progName = 
	bakeProg = path.Join(os.Getenv("BAKE"), "bin", "bake")
)

// getBinPath returns the path to the directory containing the program for this
// project
func getBinPath() (string, error) {{
	gopath := os.Getenv("GOPATH")
	if gopath == "" {{
		return fmt.Errorf("GOPATH environment variable is empty")
	}}

	for _, p := range strings.Split(gopath, os.PathSeparator) {{
		if strings.EndsWith(p, "{ProjectName}") {{
			return path.Join(p, "bin")
		}}
	}}

	return fmt.Errorf("GOPATH doesn't contain a path to {ProjectName}")
}}

func TestSuccess(t *testing.T) {{
	cmd, _, errput := runBake(t, "-n", "x", "-l", "go")

	if ! cmd.ProcessState.Success() {
		t.Fatalf("{ProjectNameLower} failed, expected success")
	}

	if len(errput) == 0 {
		t.Fatalf("Expected error, stderr was empty")
	}
}}

func run(t *testing.T, prog string, args ...string) (cmd *exec.Cmd, o, e string) {{
	cmd = exec.Command(prog, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {{
		t.Fatalf("Couldn't get %s stdout: %v", prog, err)
	}}

	stderr, err := cmd.StderrPipe()
	if err != nil {{
		t.Fatalf("Couldn't get %s stderr: %v", prog, err)
	}}

	err = cmd.Start()
	if err != nil {{
		t.Fatalf("Couldn't start %s: %v", prog, err)
	}}

	if o, err = readLines(stdout); err != nil {{
		t.Fatalf("Error reading stdout for %s: %v", prog, err)
	}}

	if e, err = readLines(stderr); err != nil {{
		t.Fatalf("Error reading stderr for %s: %v", prog, err)
	}}

	cmd.Wait()

	return
}}
{!}
