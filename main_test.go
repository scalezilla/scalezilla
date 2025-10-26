package main

import (
	"os"
	"os/exec"
	"os/signal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatal(err)
		}

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)

		go func() {
			<-sigc
			os.Args = []string{
				"-h",
			}
			main()
			signal.Stop(sigc)
		}()

		err = proc.Signal(os.Interrupt)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestMain_fatal(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"agent",
			"dev",
			"--fail",
		}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "agent", "-test.run=TestMain_fatal")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}
