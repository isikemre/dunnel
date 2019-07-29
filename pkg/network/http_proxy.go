package network

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	bodyBytes, err := ioutil.ReadAll(conn.Response.Body)
	if err != nil {
		log.Fatalln(err)
		panic("Error wile reading body")
	}
	fmt.Println(string(bodyBytes))
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

// for DEV mode
var dunnelSession = newDunnelSession("123456")

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request

	if isWebSocket(r) {
		handleWebSocket(wr, r, &dunnelSession)
		return
	}

	unixClient := unixSocketClient("/var/run/docker.sock")

	log.Printf("%v %v", r.Method, r.RequestURI)
	req, err = http.NewRequest(r.Method, "http+unix://docker"+r.RequestURI, r.Body)
	for name, value := range r.Header {
		req.Header.Set(name, value[0])
	}
	resp, err = unixClient.Do(req)

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

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	conn := &HttpConnection{r, resp}
	PrintHTTP(conn)

	r.Body.Close()
	resp.Body.Close()
}

func StartHTTPProxy(port int) {
	proxy := NewProxy()
	fmt.Println("==============================")
	err := http.ListenAndServe(":"+strconv.Itoa(port), proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
