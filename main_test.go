package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		os.Args = []string{
			"scalezilla",
			"--help",
		}
		main()
	})

	t.Run("fatal", func(t *testing.T) {
		if os.Getenv("FATAL_AGENT_DEV") == "1" {
			os.Args = []string{
				"scalezilla",
				"agent",
				"dev",
				"--fail",
			}
			main()
			return
		}

		bin, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}

		cmd := exec.Command(bin, "-test.run=TestMain/fatal")
		cmd.Env = append(os.Environ(), "FATAL_AGENT_DEV=1")
		err = cmd.Run()
		if err == nil {
			t.Fatal("expected non-nil error from fatal exit")
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() != 1 {
				t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
			}
			// success
			return
		}

		t.Fatalf("unexpected error type: %T, %v", err, err)
	})
}
