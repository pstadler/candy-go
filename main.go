package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
)

var (
	proxy_url  *url.URL
	index_html string
)

func main() {
	index_html = parseIndexFile("public/index.html")
	proxy_url, _ = url.Parse(config("HTTP_Bind"))

	http.HandleFunc("/", serveIndexFunc)
	http.HandleFunc("/http-bind/", serveProxyFunc)
	http.Handle("/candy/", http.StripPrefix("/candy/", http.FileServer(http.Dir("./public/candy/"))))

	http.ListenAndServe(config("App"), nil)
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
	var buf bytes.Buffer
	t, _ := template.ParseFiles(path)
	t.Execute(&buf, map[string]string{"core": config("Core"), "view": config("View"), "connect": config("Connect")})
	return buf.String()
}
