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
	confFile := flag.String("logconfig", "golog.yml", "Golog configuration file in yaml format")
	file := flag.String("file", "", "Log file to view")
	flag.Parse()

	//Parse yaml config file
	yamlFile, err := ioutil.ReadFile(*confFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var configuration logreader.LogConfig
	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if *file == "" {
		file = &configuration.LogFile
	}

	logReader := logreader.NewLogReader(*file, configuration)
	logDisplay := logdisplay.NewLogDisplay(&logReader)
	logDisplay.DisplayUI()
}
