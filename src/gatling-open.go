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
var successCount, totalPush int

func gatling(url string, object string, xHeaders string, ch chan<- bool, rtype string, otype string) {

	count := 0
	out := []byte(object)
	cheads := strings.Split(xHeaders, ",")

	for {

		switch strings.ToUpper(rtype) {
		case "POST":
			request, _ := http.NewRequest(strings.ToUpper(rtype), url, bytes.NewBuffer(out))
			switch otype {
			case "xml":
				request.Header.Set("Content-Type", "text/xml")
			case "json":
				request.Header.Set("Content-Type", "application/json")
			default:
				log.Printf("[*] Invalid content type. Please specify json or xml")
				os.Exit(1)
			}
			if len(cheads) > 0 {
				for _, v := range cheads {
					hs := strings.Split(v, ":")
					request.Header.Set(hs[0], hs[1])
				}
			}
			response, err := client.Do(request)
			if err != nil {
				log.Fatalf("[*] Error making POST request: %s", err)
			}
			defer response.Body.Close()

			log.Printf("[*] Request complete! Finished Request No: #%v, Status: %v", count, response.StatusCode)

			totalPush++
			count++

			if response.StatusCode == 200 || response.StatusCode == 201 {
				successCount++
			}
			ch <- true

		case "GET":
			request, _ := http.NewRequest(strings.ToUpper(rtype), url, nil)
			if len(cheads) > 0 {
				for _, v := range cheads {
					hs := strings.Split(v, ":")
					request.Header.Set(hs[0], hs[1])
				}
			}
			response, err := client.Do(request)
			if err != nil {
				log.Fatalf("[*] Error making GET request: %s", err)
			}
			defer response.Body.Close()
			log.Printf("[*] Request complete! Finished Request No: #%v, Status: %v", count, response.StatusCode)
			totalPush++
			count++

			if response.StatusCode == 200 || response.StatusCode == 201 {
				successCount++
			}
			ch <- true

		default:
			log.Printf("[*] Bad HTTP method specified. Please specify either 'GET' or 'POST' as 'rtype'")
		}

	}

}

func main() {

	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)

	url := flag.String("url", "", "Url to stress test e.g. 'http://acme.com'.")
	request := flag.Int("rps", 0, "Number of requests to make simultaneously.")
	data := flag.String("data", "", "Custom object to post e.g. {'foo':'bar'}.")
	dataType := flag.String("data-type", "", "Type of object to post. e.g. 'xml' or 'json'.")
	numRequests := flag.Int("total-requests", 0, "Total number of requests to make.")
	requestType := flag.String("type", "", "HTTP request type you'd like to make. Either 'GET' or 'POST'.")
	heads := flag.String("headers", "", "Set HTTP headers. Format should be for example 'Auth:SomeToken,X-Header:Sugar'. Headers should be separated by commas.")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\nParameters:\n", os.Args[0])
		flag.PrintDefaults()
	}

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
	switch *url {
	case "":
		flag.Usage()
	default:
		ch := make(chan bool)
		for i := 0; i < *request; i++ {
			go gatling(*url, *data, *heads, ch, *requestType, *dataType)
		}

		for r := 0; r < *numRequests; r++ {
			<-ch
		}

		elapsedTime := time.Since(start)
		failedCount := totalPush - successCount
		log.Printf("[*] Total number of successful requests: %v", successCount)
		log.Printf("[*] Total number of failed requests: %v", failedCount)
		log.Printf("[*] Total number of requests: %v", totalPush)
		log.Printf("[*] Total time elapsed: %v", elapsedTime)
	}

}
