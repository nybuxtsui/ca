package tests

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"testing"
)

const (
	binPath    = "../bin/ca"
	depotDir   = "certs-test"
	hostname   = "host1"
	passphrase = "123456"
)

func run(command string, args ...string) (string, string, error) {
	var stdoutBytes, stderrBytes bytes.Buffer
	args = append([]string{"--depot-path", depotDir}, args...)
	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	err := cmd.Run()
	return stdoutBytes.String(), stderrBytes.String(), err
}

func runWithStdin(stdin io.Reader, command string, args ...string) (string, string, error) {
	var stdoutBytes, stderrBytes bytes.Buffer
	args = append([]string{"--depot-path", depotDir}, args...)
	cmd := exec.Command(command, args...)
	cmd.Stdin = stdin
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	err := cmd.Run()
	return stdoutBytes.String(), stderrBytes.String(), err
}

func TestVersion(t *testing.T) {
	stdout, stderr, err := run(binPath, "--version")
	if stderr != "" || err != nil {
		t.Fatalf("Received unexpected error: %v, %v", stderr, err)
	}
	if !strings.Contains(stdout, "version") {
		t.Fatalf("Received unexpected stdout: %v", stdout)
	}
}
