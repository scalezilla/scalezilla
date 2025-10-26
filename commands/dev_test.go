package commands

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandsDev(t *testing.T) {
	assert := assert.New(t)

	t.Run("server_success", func(t *testing.T) {
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatal(err)
		}

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)

		cmd := Dev()

		go func() {
			<-sigc
			assert.Nil(cmd.Run(context.Background(), []string{}))
			signal.Stop(sigc)
		}()

		err = proc.Signal(os.Interrupt)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("fail", func(t *testing.T) {
		cmd := Dev()
		assert.Error(cmd.Run(context.Background(), []string{"dev", "--fail"}))
	})
}
