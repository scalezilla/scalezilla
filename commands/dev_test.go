package commands

import (
	"context"
	"os"
	"os/signal"
	"testing"
	"time"

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
		// SLEEP MUST BE KEPT to avoid data race during tests
		// it's only required with proc, err := os.FindProcess(os.Getpid())
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("fail", func(t *testing.T) {
		cmd := Dev()
		assert.Error(cmd.Run(context.Background(), []string{"dev", "--fail"}))
	})
}
