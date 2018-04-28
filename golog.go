package main

import (
	"flag"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/oskanaan/golog/logreader"
	"github.com/oskanaan/golog/logdisplay"
	"os"
	"regexp"
	"strings"
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

	if configuration.AutoDetectFiles == "Y" {
		files := populateFilePaths()
		for _, file := range files {
			configuration.Files = append(configuration.Files, logreader.LogFile{LogFile: file, Name: file[strings.LastIndex(file, `/`)+1:]})
		}
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

	if configuration.AutoDetectFiles == "Y" {
		files := populateFilePaths()
		for _, file := range files {
			configuration.Files = append(configuration.Files, logdisplay.LogFile{LogFile: file, Name: file[strings.LastIndex(file, `/`)+1:]})
		}
	}

	return configuration
}

func populateFilePaths() []string{
	var filePaths []string
	var populateFromSubDirectories func (currentDirectory string)
	populateFromSubDirectories = func (currentDirectory string) {

		files, err := ioutil.ReadDir(currentDirectory)
		if err != nil {
			log.Printf("Could not get information about working directory, error returned was %v", err)
		}

		for _, fileInfo := range files {
			if fileInfo.IsDir() {
				if fileInfo.Name()[0:1] == "." {
					continue
				}
				wd := currentDirectory+`/`+fileInfo.Name()
				populateFromSubDirectories(wd)
			} else {
				if found, _ := regexp.MatchString(`.*\.log$`, fileInfo.Name()); found {
					filePaths = append(filePaths, currentDirectory+`/`+fileInfo.Name())
				}
			}
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Couldnt get current working directory information, error returned was %v", err)
	}

	populateFromSubDirectories(dir)

	return filePaths
}