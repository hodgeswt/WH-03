package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hodgeswt/WH-03/internal/wh03"
	"github.com/hodgeswt/utilw/pkg/logw"
)

func main() {
	// TODO: all logic
	content, err := os.ReadFile("config.json")
	defaultConfig := &wh03.Config{
		ClockFreq: 2,
		RamK:      32,
		LogLevel:  "ERROR",
		LogFile:   "",
		RomK:      32,
		RomFile:   "rom.bin",
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

	rom := make([]int64, cfg.RomK*(2^10))

	file, err := os.Create(cfg.RomFile)
	if err != nil {
		panic(fmt.Sprintf("Unable to create ROM file at %s, err: %v", cfg.RomFile, err))
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, rom)

	if err != nil {
		panic(fmt.Sprintf("Unable to write ROM file at %s, err: %v", cfg.RomFile, err))
	}
}
