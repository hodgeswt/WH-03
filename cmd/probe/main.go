package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s: %s\n", msg.Topic(), msg.Payload())
}

type Config struct {
	Probe []string `json:"probe"`
}

func main() {

	content, err := os.ReadFile("probeConfig.json")

	var cfg *Config
	if err != nil {
		panic(err)
	} else {
		loadedCfg := new(Config)
		err := json.Unmarshal(content, &loadedCfg)
		if err != nil {
			panic(err)
		} else {
			cfg = loadedCfg
		}
	}

	if len(cfg.Probe) == 0 {
		fmt.Println("nothing to probe, exiting.")
		os.Exit(0)
	}

	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, topic := range cfg.Probe {
		if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}

	select {}
}
