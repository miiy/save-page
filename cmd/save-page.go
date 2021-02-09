package main

import (
	"flag"
	"fmt"
	"github.com/miiy/save-page/pkg/config"
	"github.com/miiy/save-page/pkg/file"
	"github.com/miiy/save-page/pkg/page"
	"log"
	"os"
	"os/exec"
)

func main() {
	var argT = flag.String("t", "", "-t: html or pdf")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("save-file: save-file -t html http://test.com/test")
		return
	}

	tp := *argT
	if tp != "html" && tp != "pdf" {
		fmt.Println("save-file: save-file -t pdf http://test.com/test test.pdf")
		return
	}

	var url string
	var fileName string
	for i, v := range args {
		if i == 0 {
			url = v
		}
		if i == 1 {
			fileName = args[1]
		}
	}

	cfg, err := config.NewConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	initialization(cfg)

	savePageHandler(cfg, tp, url, fileName)
}

func initialization(config *config.Config)  {
	if !file.Exists(config.StoragePath) {
		if err := file.Mkdir(config.StoragePath); err != nil {
			log.Fatal(err)
		}
	}
}

func savePageHandler(cfg *config.Config, tp, url, fileName string) {
	if tp == "html" {
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

	if tp == "pdf" {
		f := cfg.StoragePath + string(os.PathSeparator) + fileName
		cmd := exec.Command( "wkhtmltopdf", url, f)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}

}