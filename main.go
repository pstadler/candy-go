package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
)

var (
	config     *Configuration
	proxy_url  *url.URL
	index_html string
)

func main() {
	config = loadConfig("config.json")

	index_html = parseIndexFile("public/index.html")
	proxy_url, _ = url.Parse(fmt.Sprintf("%s:%d%s",
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
	r.URL = proxy_url
	r.Host = proxy_url.Host

	client := &http.Client{}
	response, _ := client.Do(r)

	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.Header().Set("Content-Encoding", response.Header.Get("Content-Encoding"))
	io.Copy(w, response.Body)
}

func parseIndexFile(path string) string {
	core, _ := json.Marshal(config.Candy.Core)
	view, _ := json.Marshal(config.Candy.View)
	connect, _ := json.Marshal(config.Candy.Connect)
	connect = connect[1 : len(connect)-1]

	var buf bytes.Buffer
	t, _ := template.ParseFiles(path)
	t.Execute(&buf, map[string]string{"core": string(core), "view": string(view), "connect": string(connect)})
	return buf.String()
}
