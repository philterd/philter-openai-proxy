package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	var o OpenAIRequest
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := 0; i < len(o.Messages); i++ {
		fmt.Fprintf(w, "Content: %+v", o.Messages[i].Content)
		// TODO: Call Philter with the content and replace the content with Philter's response.
		o.Messages[i].Content = "Who won the World Series in 1999?"
	}

	j, err := json.Marshal(o)
	new_body_content := string(j[:])

	r.Body = ioutil.NopCloser(strings.NewReader(new_body_content))

	r.ContentLength = int64(len(new_body_content))
	r.Host = p.target.Host

	p.proxy.ServeHTTP(w, r)
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {

	target, err := url.Parse("https://api.openai.com")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	p := &Proxy{target: target, proxy: proxy}

	port := getEnv("PHILTER_PROXY_PORT", "8080")
	cert_file := getEnv("PHILTER_PROXY_CERT_FILE", "cert.pem")
	key_file := getEnv("PHILTER_PROXY_KEY_FILE", "key.pem")

	err = http.ListenAndServeTLS(":"+port, cert_file, key_file, p)
	if err != nil {
		panic(err)
	}

}
