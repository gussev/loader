package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Version            string  `json:"version"`
	Profile            bool    `json:"profile"`
	Print_response     bool    `json:"print_response"`
	Files              []File `json:"files"`
	Loops              int     `json:"loops"`
	Wait               int     `json:"wait"`
	Goroutines         int     `json:"goroutines"`
	Wait_duration_type string  `json:"wait_duration_type"`
	Sleep_before_start int     `json:"sleep_in_seconds_before_start"`
}
type File struct{
	Name string `json:"file"`
}
func NewConfig(path_to_file string)(config *Config,err error){
	jsonFile,err := os.Open(path_to_file)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	json.Unmarshal(jsonData, &config)
	return
}
func (c *Config)Sleep(){
	switch c.Wait_duration_type {
		case "Microsecond": time.Sleep(time.Microsecond*time.Duration(c.Wait))
		case "Millisecond": time.Sleep(time.Millisecond*time.Duration(c.Wait))
		case "Second":      time.Sleep(time.Second*time.Duration(c.Wait))
		default: panic("wait duration type:"+c.Wait_duration_type+" is not supported")
	}
}

func (c *Config) SleepBeforeStart(){
	time.Sleep(time.Second * time.Duration(c.Sleep_before_start))
}
func (c *Config) RunProfileServer(){
	if c.Profile == false {
		return
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
func (c *Config) Print_Response(body []byte){
	if c.Print_response == false {
		return
	}
	fmt.Println(string(body))
}