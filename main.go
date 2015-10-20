package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	// Go YAML parser
	"gopkg.in/yaml.v2"
)

// Config : loaded at startup from yaml config files
type Config struct {
	NumTaps  int      `yaml:"num_taps"`
	TapNames []string `yaml:"tap_names"`
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

func apiHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/api/v1/"):]
	color := r.FormValue("color")
	fmt.Println(target, color)
	updateColor(target, color, config)
	fmt.Fprintf(w, "ok")
}

func main() {
	loadConf(&config, "settings_default.yaml", "settings_user.yaml")
	fmt.Printf("%+v\n", config)

	setupMQTT(config)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(
		http.Dir("./templates/"))))
	http.HandleFunc("/api/v1/", apiHandler)

	var listen = ":8081"
	http.ListenAndServe(listen, nil)
}
