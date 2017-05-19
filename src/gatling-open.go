package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var client *http.Client
var successful_count int
var total_push int

func gatlingXML(url string, xml string, cheads []string, ch chan<- bool, rtype string) {

	count := 0
	out := []byte(xml)

	for {
		request, _ := http.NewRequest(rtype, url, bytes.NewBuffer(out))
		for _, v := range cheads {
			hs := strings.Split(v, ":")
			request.Header.Set(hs[0], hs[1])
		}
		request.Header.Set("Content-Type", "text/xml; charset=utf-8")
		response, yoo := client.Do(request)
		if yoo != nil {
			log.Printf("[*] Error making request.")
		}
		defer response.Body.Close()

		log.Printf("[*] Request complete! Finished Request No: #%v, Status: %v", count, response.StatusCode)

		total_push += 1

		count += 1

		if response.StatusCode == 200 {
			successful_count += 1
			ch <- true
		}

	}

}

func gatlingJSON(url string, jsonObject string, custHeaders []string, ch chan<- bool, rtype string) {

	count := 0
	out := []byte(jsonObject)

	for {
		request, _ := http.NewRequest(rtype, url, bytes.NewBuffer(out))
		for _, headerVal := range custHeaders {
			s := strings.Split(headerVal, ":")
			request.Header.Set(s[0], s[1])
		}
		request.Header.Set("Content-Type", "application/json")
		response, yoo := client.Do(request)
		if yoo != nil {
			log.Printf("[*] Error making request.")
		}
		defer response.Body.Close()
		log.Printf("[*] Request complete! Finished Request No: #%v, Status: %v", count, response.StatusCode)

		total_push += 1

		count += 1

		if response.StatusCode == 201 {
			successful_count += 1
			ch <- true
		}

	}

}

func main() {

	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	urlString := flag.String("url", "", "Url to stress test")
	requestInt := flag.Int("rps", 0, "Number of requests to make simultaneously")
	requestObject := flag.String("object", "", "Custom object to post")
	//requestToken := flag.String("token", "", "Token to authorize requests")
	objType := flag.String("objectType", "", "Type of object to post")
	numRequests := flag.Int("numR", 0, "Total number of requests to make")
	reqType := flag.String("type", "", "HTTP request type you'd like to make. Currently supported type(s) 'POST'")
	heads := flag.Args()
	flag.Parse()

	// Customizing Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, yoo := defaultRoundTripper.(*http.Transport)
	if !yoo {
		panic(fmt.Sprintf("defaultRoundTripper lied! Not an *http.Transport"))
	}

	defaultTransport := *defaultTransportPointer
	defaultTransport.MaxIdleConns = 1500
	defaultTransport.MaxIdleConnsPerHost = 1500

	client = &http.Client{Transport: &defaultTransport}

	start := time.Now()
	if *objType == "json" {
		ch := make(chan bool)
		for i := 0; i < *requestInt; i++ {
			go gatlingJSON(*urlString, *requestObject, heads, ch, *reqType)
		}

		for r := 0; r < *numRequests; r++ {
			<-ch
		}
	}

	if *objType == "xml" {
		ch := make(chan bool)
		for i := 0; i < *requestInt; i++ {
			go gatlingXML(*urlString, *requestObject, heads, ch, *reqType)
		}

		for c := 0; c < *numRequests; c++ {
			<-ch
		}
	}

	elapsedTime := time.Since(start)
	var failed_count int
	failed_count = total_push - successful_count
	log.Printf("[*] Total number of successful requests: %v", successful_count)
	log.Printf("[*] Total number of failed requests: %v", failed_count)
	log.Printf("[*] Total number of requests: %v", total_push)
	log.Printf("[*] Total time elapsed: %v", elapsedTime)
}
