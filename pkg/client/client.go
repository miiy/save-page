package client

import (
	"crypto/tls"
	"github.com/miiy/save-page/pkg/config"
	"golang.org/x/net/publicsuffix"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	urlpkg "net/url"
	"path"
	"time"
)

type Client struct {
	config *config.Config
	client *http.Client
}

func NewClient(config *config.Config) (*Client, error) {
	client := &http.Client{}
	// timeout
	if config.Timeout > 0 {
		client.Timeout = time.Duration(config.Timeout) * time.Second
	}
	// cookie
	// jar, err := cookiejar.New(nil)
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client.Jar = jar
	var transport *http.Transport
	// proxy
	if config.Proxy != "" {
		proxyUrl := func(_ *http.Request) (*urlpkg.URL, error) {
			return urlpkg.Parse(config.Proxy)
		}
		if transport == nil {
			transport = &http.Transport{}
		}
		transport.Proxy = proxyUrl
	}
	//
	if config.DialContext.Timeout > 0 && config.DialContext.KeepAlive > 0 {
		if transport == nil {
			transport = &http.Transport{}
		}
		transport.DialContext = (&net.Dialer{
			Timeout:   time.Duration(config.DialContext.Timeout) * time.Second,
			KeepAlive: time.Duration(config.DialContext.KeepAlive) * time.Second,
		}).DialContext
	}

	if true {
		if transport == nil {
			transport = &http.Transport{}
		}
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if transport != nil {
		client.Transport = transport
	}
	return &Client{
		config: config,
		client: client,
	}, nil
}

func (c *Client) Get(url string, params map[string]string, headers map[string] string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			if val != "" {
				q.Add(key, val)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	if c.config.Debug == true {
		u, _ := urlpkg.Parse(url)
		log.Printf("Request url: %s headers: %s cookies: %s", url, headers, c.client.Jar.Cookies(u))
	}

	return c.client.Do(req)
}


func (c *Client) Post(url string, body io.Reader, headers map[string] string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	return c.client.Do(req)
}

func Headers(url *urlpkg.URL, headers map[string]string) (map[string]string, error) {
	var defaultHeaders = map[string]string{
		"Host":            url.Host,
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
		"Accept":          headerAccept(url),
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding": "gzip, deflate, br",
		"Referer":         "",
		"Connection":      "keep-alive",
		// "Cookie": "",
		"Upgrade-Insecure-Requests": "1",
		"Pragma" : "no-cache",
		"Cache-Control": "no-cache",
	}

	if headers != nil {
		for k, v := range headers {
			defaultHeaders[k] = v
		}
	}
	return defaultHeaders, nil
}

func headerAccept(url *urlpkg.URL) string {
	defaultAccept := "*/*"
	htmlAccept := "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	imageAccept := "image/webp,*/*"
	cssAccept := "text/css,*/*;q=0.1"
	jsAccept := defaultAccept
	extAcceptMap := map[string]string{
		".gif": imageAccept,
		".png": imageAccept,
		".jpg": imageAccept,
		".jpeg": imageAccept,
		".webp": imageAccept,
		".css": cssAccept,
		".js": jsAccept,
	}
	ext := path.Ext(url.String())
	for k, v := range extAcceptMap {
		if ext == k {
			return v
		}
	}
	return htmlAccept
}
