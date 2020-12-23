package page

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/miiy/save-page/client"
	"github.com/miiy/save-page/config"
	"github.com/miiy/save-page/file"
	"io/ioutil"
	"log"
	"net/http"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const PathSeparator = string(os.PathSeparator)

type Page struct {
	config *config.Config
	url *urlpkg.URL
	client *client.Client
}

func NewPage(config *config.Config, url string) (*Page, error) {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Page{
		config: config,
		url: u,
		client: c,
	}, nil
}

func (p *Page) Document() (*goquery.Document, error) {

	headers, err := client.Headers(p.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := request(p.client, p.url.String(), headers)
	if resp == nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				log.Println(err)
			}
		}
	}()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	return doc, err
}

func request(client *client.Client, url string, headers map[string]string) (*http.Response, error) {
	resp, err := client.Get(url, nil, headers)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := strconv.Itoa(resp.StatusCode) + " " + url
		return nil, errors.New(err)
	}
	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return resp, err
}

func (p *Page ) SavePage (doc *goquery.Document) (string, error) {
	title := doc.Find("head title").First().Text()
	title = file.SafeName(title)
	if title == "" {
		return "", errors.New("page title is empty")
	}

	fileName := title + ".htm"
	f := strings.TrimSuffix(p.config.StoragePath, PathSeparator) + PathSeparator + fileName

	docHtml, err := doc.Html()
	if err != nil {
		return "", err
	}
	if docType := htmlDocType(docHtml); docType != "" {
		docHtml = strings.Replace(docHtml, docType, docType + "\n<!-- saved from save-page url=" + p.url.String() + " -->", 1)
	}

	if !saveFile(f, []byte(docHtml)) {
		return "", errors.New("page save error")
	}

	return f, nil
}

func (p *Page) SaveResource(f string) {

	htmlByte, err := file.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	html := string(htmlByte)

	fileName := filepath.Base(f)
	dir := filepath.Dir(f)

	title := strings.TrimSuffix(fileName, ".htm")

	buf := bytes.NewBuffer(htmlByte)
	doc, err := goquery.NewDocumentFromReader(buf)
	if doc == nil {
		log.Fatal("doc is nil")
	}

	sourceUrl := htmlSourceUrl(html)
	if sourceUrl == "" {
		log.Fatal("No match to source url.")
	}

	// x_files
	resourceDirName := title + "_files"
	if !makeDir(dir + PathSeparator + resourceDirName) {
		log.Fatal("Make resource dir error.")
	}

	baseTagS := doc.Find("base")
	baseTagHref, ok:= baseTagS.Attr("href")
	if ok {
		baseTagS.Remove()
	}

	wg := sync.WaitGroup{}
	doc.Find("script, link, img").Each(func(i int, s *goquery.Selection) {
		//h, err := goquery.OuterHtml(s)
		//if err != nil {
		//	log.Println(err)
		//}
		//log.Println(h)

		// script src, link href, img src
		nodeName := goquery.NodeName(s)
		var attrName string
		tagSrcMap := map[string]string{
			"script": "src",
			"link": "href",
			"img": "src",
		}
		for k, v := range tagSrcMap {
			if nodeName == k {
				attrName = v
				break
			}
		}
		link, ok := s.Attr(attrName)
		if !ok {
			return
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			resourceUrl, err := parseResourceUrl(link, sourceUrl, baseTagHref)
			if err != nil {
				log.Println(err)
				return
			}

			headers, err := client.Headers(resourceUrl,  map[string]string{
				"Referer" : sourceUrl,
			})
			if err != nil {
				log.Println(err)
				return
			}
			resp, err:= request(p.client, resourceUrl.String(), headers)
			if err != nil {
				log.Println(err)
				return
			}
			defer func() {
				if resp != nil {
					if err := resp.Body.Close(); err != nil {
						log.Println(err)
					}
				}
			}()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}

			resourceFileName := filepath.Base(link)
			resourceFileName = file.SafeName(resourceFileName)
			sFile := dir + PathSeparator + resourceDirName + PathSeparator + resourceFileName
			saveFile(sFile, b)
			rSFile :=  "./" + resourceDirName + "/" + resourceFileName
			s.SetAttr(attrName, rSFile)
		}()

	})
	wg.Wait()
	if p.config.Debug {
		log.Println("Resource download complete.")
	}
	docHtml, _:= doc.Html()

	err = file.Write(f, []byte(docHtml))
	if err != nil {
		log.Println(err)
	}
	log.Println("Completed.")
}

func parseResourceUrl(urlRaw, sourceUrl, baseTagHref string) (*urlpkg.URL, error) {
	urlU, err := urlpkg.QueryUnescape(urlRaw)
	if err != nil {
		return nil, err
	}
	url, err := urlpkg.Parse(urlU)
	if err != nil {
		return nil, err
	}
	sourceUrlP, err := urlpkg.Parse(sourceUrl)
	if err != nil {
		return nil, err
	}
	if !url.IsAbs() {
		url.Scheme = sourceUrlP.Scheme
		url.Host = sourceUrlP.Host
		if baseTagHref != "" {
			url.Path = baseTagHref + strings.TrimPrefix(url.Path, baseTagHref)
		}
	}
	return url, nil
}

func makeDir(name string) bool {
	if file.Exists(name) {
		log.Println(name + " already exist.")
		return false
	}
	err := file.Mkdir(name)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func saveFile(name string, b []byte) bool {
	if file.Exists(name) {
		log.Println(name + " already exist.")
		return false
	}
	err := file.Write(name, b)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func htmlDocType(html string) string {
	expr := `<![DOCTYPE|doctype].*?>`
	re, err := regexp.Compile(expr)
	if err != nil {
		log.Println(err)
	}
	return re.FindString(html)
}

func htmlSourceUrl(html string) string {
	expr := `<!-- saved from.*?url=(.*?) -->`
	re, err := regexp.Compile(expr)
	if err != nil {
		log.Println(err)
	}
	mached := re.FindStringSubmatch(html)
	if len(mached) > 1{
		return mached[1]
	}
	return ""
}
