package main

import (
	"flag"
	"io/ioutil"
	"log"

	app "github.com/1PALADIN1/gigachat_server/internal"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "configs/server_config.yaml"

func main() {
	flag.Parse()

	configPath := defaultConfigPath
	if len(flag.Args()) != 0 {
		configPath = flag.Args()[0]
	}

	configText, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config := &app.Config{}
	yaml.Unmarshal(configText, config)

	app.Run(config)
}
