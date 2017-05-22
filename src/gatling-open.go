package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var client *http.Client
var successful_count int
var total_push int

func gatling(url string, object string, xHeaders string, ch chan<- bool, rtype string, otype string) {

	count := 0
	out := []byte(object)
	cheads := strings.Fields(xHeaders)
	request, _ := http.NewRequest(rtype, url, bytes.NewBuffer(out))
	if strings.ToLower(otype) != "xml" || strings.ToLower(otype) != "json" {
		log.Printf("[*] Invalid content type. Please specify json or xml")
		os.Exit(1)
	} else {
		request.Header.Set("Content-Type", "application/"+otype)
	}
	for _, v := range cheads {
		hs := strings.Split(v, ":")
		request.Header.Set(hs[0], hs[1])
	}
	response, yoo := client.Do(request)
	if yoo != nil {
		log.Fatalf("[*] Error making request.", yoo)
	}
	defer response.Body.Close()

	log.Printf("[*] Request complete! Finished Request No: #%v, Status: %v", count, response.StatusCode)

	total_push += 1

	count += 1

	if response.StatusCode == 200 || response.StatusCode == 201 {
		successful_count += 1
		ch <- true
	}
}

func main() {

	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	urlString := flag.String("url", "", "Url to stress test e.g. 'http://acme.com'")
	requestInt := flag.Int("rps", 0, "Number of requests to make simultaneously")
	requestObject := flag.String("object", "", "Custom object to post e.g. {'foo':'bar'}")
	objType := flag.String("objectType", "", "Type of object to post. e.g. 'xml' or 'json'")
	numRequests := flag.Int("numR", 0, "Total number of requests to make")
	reqType := flag.String("type", "", "HTTP request type you'd like to make. Currently supported type(s) 'POST'")
	heads := flag.String("headers", "", "HTTP headers to include when making request. Format should be for example 'Auth:SomeToken X-Header:Sugar'. \nHeaders should be separated by spaces")
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
	switch *objType {
	case "":
		flag.PrintDefaults()
	default:
		ch := make(chan bool)
		for i := 0; i < *requestInt; i++ {
			go gatling(*urlString, *requestObject, *heads, ch, *reqType, *objType)
		}

		for r := 0; r < *numRequests; r++ {
			<-ch
		}

		elapsedTime := time.Since(start)
		failed_count := total_push - successful_count
		log.Printf("[*] Total number of successful requests: %v", successful_count)
		log.Printf("[*] Total number of failed requests: %v", failed_count)
		log.Printf("[*] Total number of requests: %v", total_push)
		log.Printf("[*] Total time elapsed: %v", elapsedTime)
	}

}
