package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	mlog "github.com/tahersoft-go/kit/log"
)

var ignoreSignals = []os.Signal{os.Interrupt}
var forwardSignals = []os.Signal{syscall.SIGTERM}

func Start(rootCmd *cobra.Command) error {

	wd, err := os.Getwd()
	if err != nil {
		// It would be very strange to end up here
		return fmt.Errorf("failed to determine current working directory: %s", err)
	}

	// Create Meta for CMDs
	logger := mlog.New().With(rootCmd.Context())
	meta := &Meta{
		CallerContext: rootCmd.Context(),
		WorkingDir:    wd,
		ShutdownCh:    makeShutdownCh(),
		Logger:        logger,
	}

	serverCommand := &ServerCommand{meta}
	serverCommand.AddCommandToCobra(rootCmd)

	return rootCmd.Execute()
}

// makeShutdownCh creates an interrupt listener and returns a channel.
// A message will be sent on the channel for every interrupt received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, ignoreSignals...)
	signal.Notify(signalCh, forwardSignals...)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}
