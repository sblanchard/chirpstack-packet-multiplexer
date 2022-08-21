package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sblanchard/chirpstack-packet-multiplexer/internal/config"
	"github.com/sblanchard/chirpstack-packet-multiplexer/internal/multiplexer"
)

var mp *multiplexer.Multiplexer

func run(cmd *cobra.Command, args []string) error {

	tasks := []func() error{
		setLogLevel,
		printStartMessage,
		setupMultiplexer,
	}

	for _, t := range tasks {
		if err := t(); err != nil {
			log.Fatal(err)
		}
	}

	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping chirpstack-packet-multiplexer")
		if err := mp.Close(); err != nil {
			log.Fatal(err)
		}
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

	return nil
}

func setLogLevel() error {
	log.SetLevel(log.Level(uint8(config.C.General.LogLevel)))
	return nil
}

func printStartMessage() error {
	log.WithFields(log.Fields{
		"version": version,
		"docs":    "https://www.chirpstack.io/",
	}).Info("starting ChirpStack Packet Multiplexer")
	return nil
}

func setupMultiplexer() error {
	var err error
	mp, err = multiplexer.New(config.C.PacketMultiplexer)
	if err != nil {
		return errors.Wrap(err, "new multiplexer error")
	}

	return nil
}
