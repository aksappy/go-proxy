package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

var routeURL string

type configuration struct {
	RootURL    string `yaml:"root-url"`
	SigningKey string `yaml:"jwt-key"`
	Port       int    `yaml:"port"`
	RootDir    string `yaml:"root-dir"`
}

func main() {
	var conf configuration
	conf.ReadConfig("./config.yaml")
	routeURL = conf.RootURL
	port := fmt.Sprintf(":%d", conf.Port)
	r := mux.NewRouter()
	r.HandleFunc("/api/sec/data/{table}", mainHandler)
	r.HandleFunc("/api/files/{dir}/{id}", fileHandler)
	r.HandleFunc("/api/auth", authHandler)
	r.HandleFunc("/api/refresh", refreshHandler)
	http.Handle("/", r)
	http.ListenAndServe(port, nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {

}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dir := vars["dir"]
	id := vars["id"]
	w.Write([]byte(dir))
	w.Write([]byte(id))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := r.URL.RawQuery
	url := fmt.Sprintf(routeURL, vars["table"], q)
	switch r.Method {
	case "GET":
		performGet(w, url)
		break
	case "POST":
		performPost(w, url, r.Body)
		break
	case "PUT":
		performPut(w, url, r.Body)
		break
	case "DELETE":
		break
	default:

	}
}

func performGet(w http.ResponseWriter, url string) {
	res, err := http.Get(url)
	if err != nil {

	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {

	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// initialize http client
var client = &http.Client{}

func performPut(w http.ResponseWriter, url string, data io.ReadCloser) {
	d, _ := ioutil.ReadAll(data)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(d))
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {

	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func performPost(w http.ResponseWriter, url string, data io.ReadCloser) {
	d, _ := ioutil.ReadAll(data)
	res, err := http.Post(url, "application/json", bytes.NewReader(d))
	if err != nil {

	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {

	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

//ReadConfig reads the config.yaml and sets the respective configuration properties
func (conf *configuration) ReadConfig(filePath string) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("Missing configuration file")
	}
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		panic("Unreadable configuration file")
	}
}
