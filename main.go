package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	config     *Configuration
	index_html string
	proxy_target *url.URL
)

func main() {
	config = loadConfig("config.json")

	index_html = parseIndexFile("public/index.html")
	proxy_target, _ = url.Parse(fmt.Sprintf("%s:%d%s",
		config.HTTP_Bind.Host, config.HTTP_Bind.Port, config.HTTP_Bind.Path))

	http.HandleFunc("/", serveIndexFunc)
	http.HandleFunc("/http-bind/", serveProxyFunc)
	http.Handle("/candy/", http.StripPrefix("/candy/", http.FileServer(http.Dir("./public/candy/"))))

	http.ListenAndServe(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port), nil)
}

func serveIndexFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, index_html)
}

func serveProxyFunc(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""
	r.URL = proxy_target
	r.Host = proxy_target.Host

	client := &http.Client{}
	response, _ := client.Do(r)

	defer response.Body.Close()

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	io.Copy(w, reader)
}

func parseIndexFile(path string) string {
	buf, _ := ioutil.ReadFile(path)
	str := string(buf)

	// TODO: make this more pretty
	core_config, _ := json.Marshal(config.Candy.Core)
	view_config, _ := json.Marshal(config.Candy.View)
	connect_config, _ := json.Marshal(strings.Join(config.Candy.Connect, ","))
	str = strings.Replace(str, "OPTIONS", "{core:"+string(core_config)+",view:"+string(view_config)+"}", 1)
	str = strings.Replace(str, "CONNECT", string(connect_config), 1)
	return str
}
