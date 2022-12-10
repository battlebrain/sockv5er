package utils

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"os"
)

func showIntro() {
	fmt.Println(`
                         _           _____
                        | |         | ____|          
	  ___  ___   ___| | ____   _| |__   ___ _ __ 
	 / __|/ _ \ / __| |/ /\ \ / /___ \ / _ \ '__|
	 \__ \ (_) | (__|   <  \ V / ___) |  __/ |   
	 |___/\___/ \___|_|\_\  \_/ |____/ \___|_|
	`)
	fmt.Println("Downloading counties/regions list from AWS...")
	fmt.Println("Please wait a moment.")
}

func getRegionsAndCountries() []map[string]string {
	countryOptions, err := GenerateCountryOptions()
	if err != nil {
		log.Fatalf("Generating country options failed with error: %s\n", err)
	}
	if len(countryOptions) < 1 {
		log.Fatalf("Generating country options failed please try again after some time.")
	}
	return countryOptions
}

func showRegionsOptions(countryOptions []map[string]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Country", "Region"})
	for i := range countryOptions {
		t.AppendRows([]table.Row{{i + 1, countryOptions[i]["country"], countryOptions[i]["region"]}})
		t.AppendSeparator()
	}
	t.Render()
}

func getRegionFromUserInput(countryOptions []map[string]string, selection int) (string, error) {
	regionID := selection - 1
	region := countryOptions[regionID]["region"]
	return region, nil
}

func getUserInput(numberOfRegions int, in *os.File) int {
	if in == nil {
		in = os.Stdin
	}
	regionsRange := numberOfRegions
	var regionID = 0
	fmt.Println("Enter the id of the region in which you need to create the socks v5 proxy on.")
	fmt.Printf("Default is 1. Range 1-%d: ", numberOfRegions)
	for {
		_, err := fmt.Fscanf(in, "%d\n", &regionID)
		if err != nil {
			log.Fatalf("Unexpected input. %s", err)
		} else if regionID > 0 && regionID <= regionsRange {
			break
		} else {
			fmt.Printf("Please choose a number between 1 and %d: ", numberOfRegions)
		}
	}
	return regionID
}

func createSocksV5Tunnel() {
	config := SSHConfig{}
	e := ENVData{}
	settings, _ := e.Read()
	config.SocksV5IP = "127.0.0.1"
	config.SocksV5Port = "1337"
	config.SSHHost = "gmail.com"
	config.SSHPort = "22"
	config.PrivateKey = ReadFileContent(settings.PrivateKeyPath)
	config.KnownHostsFilepath = settings.SSHKnownHostsPath
	config.SSHUsername = "ubuntu"
	config.StartSocksV5Server()
}

func StartWorker() {
	showIntro()
	countryOptions := getRegionsAndCountries()
	showRegionsOptions(countryOptions)
	selection := getUserInput(len(countryOptions), nil)
	region, err := getRegionFromUserInput(countryOptions, selection)
	if err != nil {
		log.Fatalf("SockV5er failed with error: %s", err)
	}
	fmt.Printf("Selected region: %s\n", region)
	createSocksV5Tunnel()
}
