package main

import (
	"os"
	"os/signal"

	"github.com/hodgeswt/WH-03/internal/types"
	"github.com/hodgeswt/WH-03/internal/wh03"
	"github.com/hodgeswt/utilw/pkg/logw"
)

func main() {
	logw.Debug("^main.main")
	defer logw.Debug("$main.main")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	types.Broker.Init(10)

	cpu := new(wh03.CPU)
	cpu.Cfg = &wh03.Config{
		ClockFreq: 2,
	}

	go handleSigint(sigint, cpu)
	go cpu.Run()

	logw.Info("main.main: WH-03 Started")

	// Wait forever
	select {}
}

func handleSigint(sigint chan os.Signal, main *wh03.CPU) {
	_ = <-sigint

    logw.Info("main.handleSigint: sigint received")

	main.Stop()

	logw.Info("main.handleSigint: WH-03 Stopped")

	os.Exit(0)
}
