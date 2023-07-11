package main

import (
	"bytes"
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

type FilterResponse struct {
	FilteredText string `json:"filteredText"`
	Context      string `json:"context"`
	DocumentId   string `json:"documentId"`
}

func Filter(endpoint string, input string, context string, documentId string, filterProfile string) FilterResponse {

	var text = []byte(input)

	base, err := url.Parse(endpoint + "/api/filter")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	params := url.Values{}
	params.Add("c", context)
	params.Add("d", documentId)
	params.Add("p", filterProfile)

	base.RawQuery = params.Encode()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest("POST", base.String(), bytes.NewReader(text))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	request.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}
	response, err := client.Do(request)

	documentId = response.Header.Get("x-document-id")

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	return FilterResponse{FilteredText:string(responseData), Context:context, DocumentId:documentId}

}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	philter_endpoint := getEnv("PHILTER_ENDPOINT", "https://10.0.2.51:8080")
	log.Println("Proxying request to " + philter_endpoint)

	var o OpenAIRequest
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := 0; i < len(o.Messages); i++ {
		filterResponse := Filter(philter_endpoint, o.Messages[i].Content, "context", "docid", "default")
		o.Messages[i].Content = filterResponse.FilteredText
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

    fmt.Println("Started philter-openai-proxy on port " + port)
	err = http.ListenAndServeTLS(":"+port, cert_file, key_file, p)

	if err != nil {
		panic(err)
	}

}
