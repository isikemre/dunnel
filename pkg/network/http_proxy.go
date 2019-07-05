package network

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Proxy struct{}

type HttpConnection struct {
	Request  *http.Request
	Response *http.Response
}

type HttpConnectionChannel chan *HttpConnection

var connChannel = make(HttpConnectionChannel)

func NewProxy() *Proxy { return &Proxy{} }

func PrintHTTP(conn *HttpConnection) {
	fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
	for k, v := range conn.Request.Header {
		fmt.Println(k, ":", v)
	}
	fmt.Println("==============================")
	fmt.Printf("HTTP/1.1 %v\n", conn.Response.Status)
	for k, v := range conn.Response.Header {
		fmt.Println(k, ":", v)
	}
	bytes, err := ioutil.ReadAll(conn.Response.Body)
	if err != nil {
		log.Fatalln(err)
		panic("Error wile reading body")
	}
	fmt.Println(string(bytes))
	fmt.Println("==============================")
}

func HandleHTTP() {
	for {
		select {
		case conn := <-connChannel:
			PrintHTTP(conn)
		}
	}
}

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := unixSocketClient("/var/run/docker.sock")

	log.Printf("%v %v", r.Method, r.RequestURI)
	req, err = http.NewRequest(r.Method, "http+unix://docker"+r.RequestURI, r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}
	resp, err = client.Do(req)
	defer r.Body.Close()

	// combined for GET/POST
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		wr.Header().Set(k, v[0])
	}
	wr.WriteHeader(resp.StatusCode)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	io.Copy(wr, resp.Body)
	defer resp.Body.Close()

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	conn := &HttpConnection{r, resp}
	PrintHTTP(conn)
}

func StartHTTPProxy() {
	proxy := NewProxy()
	fmt.Println("==============================")
	err := http.ListenAndServe(":12345", proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
