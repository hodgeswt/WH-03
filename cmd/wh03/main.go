package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/hodgeswt/WH-03/internal/wh03"
	"github.com/hodgeswt/utilw/pkg/logw"
)

func main() {
	content, err := os.ReadFile("config.json")
	defaultConfig := &wh03.Config{
		ClockFreq:    2,
		RamK:         32,
		LogLevel:     "ERROR",
		LogFile:      "",
		RomFile:      "rom.bin",
		ProbeEnabled: false,
	}

	var cfg *wh03.Config
	if err != nil {
		logw.Warn("main.main - unable to find config.json. Using defaults.")
		cfg = defaultConfig
	} else {
		loadedCfg := new(wh03.Config)
		err := json.Unmarshal(content, &loadedCfg)
		if err != nil {
			cfg = defaultConfig
		} else {
			cfg = loadedCfg
		}
	}

	logw.SetLogLevel(cfg.LogLevel)

	if cfg.LogFile != "" {
		file, err := logw.SetOutFile(cfg.LogFile)

		if err != nil {
			panic(err)
		}

		defer file.Close()
	}

	logw.Debug("^main.main")
	defer logw.Debug("$main.main")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	wh03.Broker.Init(10, cfg.ProbeEnabled)

	cpu := new(wh03.CPU)
	cpu.Cfg = cfg

    ctx, cancel := context.WithCancel(context.Background())
	go handleSigint(sigint, cpu, cancel)
    go wh03.Broker.Run(ctx)
	go cpu.Run()

	logw.Info("main.main - WH-03 Started")

	// Wait forever
	select {}
}

func handleSigint(sigint chan os.Signal, main *wh03.CPU, cancel context.CancelFunc) {
	_ = <-sigint

	logw.Info("main.handleSigint - sigint received")

	main.Stop()
    cancel()
    wh03.Broker.Close()

	logw.Info("main.handleSigint - WH-03 Stopped")

	os.Exit(0)
}
