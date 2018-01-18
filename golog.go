package main

import (
	"flag"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/oskanaan/golog/logreader"
	"github.com/oskanaan/golog/logdisplay"
)

func main() {
	//Read command line arguments
	confFile := flag.String("logconfig", "golog.yml", "Golog logReaderConfig file in yaml format")
	flag.Parse()

	logReaderConfig := logreaderConfig(confFile)
	logDisplayConfig := logdisplayConfig(confFile)

	logReader := logreader.NewLogReader(logReaderConfig)
	logDisplay := logdisplay.NewLogDisplay(&logReader, &logDisplayConfig)
	logDisplay.DisplayUI()
}

func logreaderConfig(confFile *string) logreader.LogReaderConfig {
	yamlFile, err := ioutil.ReadFile(*confFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var configuration logreader.LogReaderConfig
	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return configuration
}

func logdisplayConfig(confFile *string) logdisplay.LogDisplayConfig {
	yamlFile, err := ioutil.ReadFile(*confFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var configuration logdisplay.LogDisplayConfig
	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return configuration
}
