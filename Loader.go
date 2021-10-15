package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"sync"
)


var cache = NewCache()
func main() {
	config,err := readConfig("config//config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
	config.SleepBeforeStart()
	config.RunProfileServer()
	wg := sync.WaitGroup{}
	for i :=0; i< config.Goroutines; i++{
		wg.Add(1)
		go run_me(config,&wg)
	}
	wg.Wait()
}

func run_me(c* Config,wg *sync.WaitGroup){
	defer wg.Done()
	for i :=0; i< c.Loops; i++{
		files := c.Files
		for _, one := range files {
			err := file_processing(one.Name,c)
			if err != nil{
				fmt.Println(err)
			}
		}
		c.Sleep()
	}
}
type convert func (name string) (str_image *string,err error)

func to_string(name string) (str_image *string,err error){
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		return nil,err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil,err
	}

	buf2 := new(bytes.Buffer)
	png.Encode(buf2, img)

	encoded := "data:image/png;base64,"+base64.StdEncoding.EncodeToString(buf2.Bytes())
	buf2.Reset()
	return &encoded,nil
}
func file_processing(name string,c* Config) (err error){
	str_image,err := cache.Get(name,to_string)
	if str_image == nil {
		return errors.New("not able to read file:"+name)
	}

	canvasdata := url.Values{}
	canvasdata.Set("canvasdata", *str_image)
	resp, err := http.PostForm("http://localhost:8080/digit", canvasdata)
	delete(canvasdata,"canvasdata")
	canvasdata = nil
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	c.Print_Response(body)
	return nil
}

func readConfig(path_to_file string) (*Config,error){
	return NewConfig(path_to_file)
}