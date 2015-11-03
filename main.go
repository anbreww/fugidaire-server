package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// Go YAML parser
	"gopkg.in/yaml.v2"
)

// Config : loaded at startup from yaml config files
type Config struct {
	NumTaps  int      `yaml:"num_taps"`
	TapNames []string `yaml:"tap_names"`
	APIKey   string   `yaml:"api_key"`
	MQTT     struct {
		Host     string
		Port     string
		Topic    string
		Username string
		Password string
		ClientID string `yaml:"client_id"`
	}
}

var config Config

func loadConf(config *Config, defaultFilename, userFilename string) {
	parseConf(defaultFilename, config) // load defaults
	// TODO: don't panic if user file doesn't exist.
	if userFilename != "" {
		parseConf(userFilename, config) // load user configuration
	}
}

func parseConf(filename string, config *Config) {
	// parse configuration file including
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, config)
	if err != nil {
		panic(err)
	}
}

// TODO : route branches between /color/ and /fan/ and /beer/
func apiHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/api/v1/"):]
	color := r.FormValue("color")
	apikey := r.FormValue("apikey")
	if apikey == config.APIKey {
		fmt.Println(target, color)
		updateColor(target, color, config)
		fmt.Fprintf(w, "ok")
	} else {
		log.Println("incorrect api key :", apikey)
		http.Error(w, "API Key missing or incorrect", http.StatusUnauthorized)
		return
	}
}

func kill(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func main() {
	loadConf(&config, "settings_default.yaml", "settings_user.yaml")
	fmt.Printf("%+v\n", config)

	setupMQTT(config)

	http.Handle("/", http.FileServer(http.Dir("./templates/")))
	http.HandleFunc("/api/v1/", apiHandler)
	http.HandleFunc("/kill", kill)

	var listen = ":8082"
	http.ListenAndServe(listen, nil)
}
