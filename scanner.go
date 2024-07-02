package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"log"
	"sync"
	"strings"
)


func RequestSecurityTxt(domain string) {
	urls := [4]string{}
	urls[0] = fmt.Sprintf("https://%s", domain)
	urls[3] = fmt.Sprintf("http://%s", domain)

	for _,requestURL := range urls {
		res, err := http.Get(requestURL)
		if err != nil {
			log.Printf("error making http request to %s: %s\n", domain, err)
			continue
		}


		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("client: could not read response body: %s\n", err)
			continue
		}


		if res.StatusCode == 200 && strings.Contains(string(resBody), "cdn.polyfill.io") {
			log.Printf("cdn.polyfill.io found on: %s\n", requestURL)
			fileName := fmt.Sprintf("raw/%s.txt", domain)
			os.WriteFile(fileName, resBody, 0644)
			return
		}
	}
	

}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing argument for input file of domains. Usage:\n%s <domains to scan>\n", os.Args[0])
		os.Exit(1)
	}

	var wg sync.WaitGroup

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for _, domain := range strings.Split(string(data[:]), "\n") {

		fileName := fmt.Sprintf("raw/%s.txt", domain)

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			wg.Add(1) // increment wait group counter
			go func (d string) {
				defer wg.Done()
				RequestSecurityTxt(d)
			} (domain)
		}

		
	}

	// wait for all groups to finish
	wg.Wait()

}