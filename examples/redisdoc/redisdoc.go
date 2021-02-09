package redisdoc

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/miiy/save-page/pkg/file"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

type link struct {
	title string
	href string
}

type RedisDoc struct {

}

func NewRedisDoc() *RedisDoc {
	return &RedisDoc{}
}

const BaseUrl = "http://redisdoc.com"

func (l *RedisDoc) Get () error {
	res, err := http.Get(BaseUrl + "/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var links []link
	doc.Find(".toctree-wrapper>ul").Children().Each(func(i int, s *goquery.Selection) {
		s.Each(func(ci int, cs *goquery.Selection) {
			linkS := s.Find("a").First()
			pre := strconv.Itoa(i + 1)
			title := linkS.Text()
			title = file.SafeName(title)
			href, _ := linkS.Attr("href")
			if err != nil {
				log.Fatal(err)
			}
			links = append(links, link{
				title: pre + "-" + title,
				href: href,
			})

			subUl := s.Has("ul")
			if subUl == nil {
				return
			}
			s.Find("ul>li").Each(func(si int, ss *goquery.Selection) {
				sLink := ss.Find("a").First()
				sPre := strconv.Itoa(si + 1)
				sTitle := sLink.Text()
				sTitle = file.SafeName(sTitle)
				sHref, _ := sLink.Attr("href")
				links = append(links, link{
					title: pre + "." + sPre + "-" + sTitle,
					href: sHref,
				})
			})
		})
	})

	for _, v := range links {
		v.href = BaseUrl + "/" + v.href
		fmt.Printf("%s %s\n", v.title, v.href)
	}
	fmt.Printf("Total: %d\n", len(links))

	fmt.Println("Begin download.")
	wg := sync.WaitGroup{}
	ch := make(chan int, 4)
	for _, v := range links {
		v.href = BaseUrl + "/" + v.href

		file := "/app/data/" + v.title + ".pdf"
		if _, err := os.Stat(file); err == nil {
			continue
		}
		log.Printf("%s %s\n", v.title, v.href)
		wg.Add(1)
		ch <- 1
		go func(v *link) {
			defer wg.Done()

			cmd := exec.Command("wkhtmltopdf", v.href, file)
			log.Println(cmd.String())
			err = cmd.Run()
			if err != nil {
				log.Println(cmd.Stderr)
				log.Println(err)
			}

			<- ch
		}(&v)
	}
	wg.Wait()
	fmt.Println("Finished.")
	return nil
}
