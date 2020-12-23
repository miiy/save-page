package main

import (
	"fmt"
	"github.com/miiy/save-page/config"
	"github.com/miiy/save-page/file"
	"github.com/miiy/save-page/page"
	"log"
	"os"
)

func main()  {
	if len(os.Args) < 2 {
		fmt.Println("save-file: try save-file http://test.com/test")
		return
	}
	url := os.Args[1]

	cfg, err := config.NewConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	initialization(cfg)

	p, err := page.NewPage(cfg, url)
	if err != nil {
		log.Fatal(err)
	}
	doc, err:= p.Document()
	if err != nil {
		log.Fatal(err)
	}
	f, err := p.SavePage(doc)
	if err != nil {
		log.Fatal(err)
	}
	p.SaveResource(f)
}

func initialization(config *config.Config)  {
	if !file.Exists(config.StoragePath) {
		if err := file.Mkdir(config.StoragePath); err != nil {
			log.Fatal(err)
		}
	}
}