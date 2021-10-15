package main

import (
	"bytes"
	"encoding/base64"
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



func main() {
	config,err := readConfig("config//config.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
	config.SleepBeforeStart()
	config.RunProfileServer()
	cache := NewCache()
	run_me := sync.NewCond(&sync.Mutex{})
	run_me.L.Lock()
	prepare_images(config,cache,run_me)
	fmt.Println("about to wait")
	run_me.Wait()
	fmt.Println("leaving wait")
	wg := sync.WaitGroup{}

	run_loader(cache ,config, &wg)
	wg.Wait()
}

func run_loader(cache *Cache,config* Config,wg *sync.WaitGroup){
	for i :=0; i< config.Goroutines; i++{
		wg.Add(1)
		go run_image(cache,config,wg)
	}
}
func run_image(cache *Cache,config* Config,wg *sync.WaitGroup){
	defer wg.Done()
	for i :=0; i< config.Loops; i++{
		walk_over_images(cache ,config)
		config.Sleep()
	}
}

func walk_over_images(cache *Cache, config *Config) {
	for name := range cache.Data(){
		mfile := cache.GetMFile(name)
		if mfile.str_image == nil{
			fmt.Println("name:",name," not ready yet")
			continue
		}
		canvasdata := url.Values{}
		canvasdata.Set("canvasdata", *(mfile.str_image))
		resp, err := http.PostForm("http://localhost:8080/digit", canvasdata)
		delete(canvasdata,"canvasdata")
		canvasdata = nil
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		config.Print_Response(body)
	}
}

func prepare_images(c* Config,cache *Cache,run_me *sync.Cond){
	files := c.Files
	for _, one := range files {
		go file_processing(one.Name,cache,run_me)
	}
}

func file_processing(name string,cache * Cache,run_me *sync.Cond) {
	mfile := cache.GetMFile(name)
	mfile.mx.Lock()
	defer mfile.mx.Unlock()
	defer run_me.Signal()
	mfile.str_image, _ = to_string(name)
}
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

func readConfig(path_to_file string) (*Config,error){
	return NewConfig(path_to_file)
}