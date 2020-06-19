package model

/*
Adaptation of https://github.com/tomnomnom/httprobe
Credits to tomnomnom
*/

import (
	"bufio"
	// "bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"

	// "log"
	"net"
	"net/http"
	"os"

	// "os/exec"
	// "strconv"
	"strings"
	"sync"
	"time"
)

func HttProbe(s *Scan, args ...interface{}) (interface{}, string, error) {
	// var result []string
	// var raw string

	config := map[string]interface{}{
		"concurrency": s.Threads,
		"skipDefault": false,
		"to":          10000,
		"preferHTTPS": false,
		"method":      "GET",
		"domains":     s.Scope,
		// "probes":      "large",
		"probes": s.Options["type"],
	}

	res := run(config)

	return map[string]interface{}{
		"result": res,
	}, "", nil
}

func run(params interface{}) []string {

	var result []string
	p := params.(interface{}).(map[string]interface{})

	domains := p["domains"].([]string)

	concurrency := p["concurrency"].(int)

	// probe flags
	probes := [1]string{p["probes"].(string)}

	// skip default probes flag
	skipDefault := p["skipDefault"].(bool)

	// timeout flag
	to := p["to"].(int)

	// prefer https
	preferHTTPS := p["preferHTTPS"].(bool)

	// HTTP method to use
	method := p["method"].(string)
	// make an actual time.Duration out of the timeout
	timeout := time.Duration(to * 1000000)

	var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: re,
		Timeout:       timeout,
	}

	// domain/port pairs are initially sent on the httpsURLs channel.
	// If they are listening and the --prefer-https flag is set then
	// no HTTP check is performed; otherwise they're put onto the httpURLs
	// channel for an HTTP check.
	httpsURLs := make(chan string)
	httpURLs := make(chan string)
	output := make(chan string)

	// HTTPS workers
	var httpsWG sync.WaitGroup
	for i := 0; i < concurrency/2; i++ {
		httpsWG.Add(1)

		go func() {
			for url := range httpsURLs {

				// always try HTTPS first
				withProto := "https://" + url
				if isListening(client, withProto, method) {
					output <- withProto

					// skip trying HTTP if --prefer-https is set
					if preferHTTPS {
						continue
					}
				}

				httpURLs <- url
			}

			httpsWG.Done()
		}()
	}

	// HTTP workers
	var httpWG sync.WaitGroup
	for i := 0; i < concurrency/2; i++ {
		httpWG.Add(1)

		go func() {
			for url := range httpURLs {
				withProto := "http://" + url
				if isListening(client, withProto, method) {
					output <- withProto
					continue
				}
			}

			httpWG.Done()
		}()
	}

	// Close the httpURLs channel when the HTTPS workers are done
	go func() {
		httpsWG.Wait()
		close(httpURLs)
	}()

	// Output worker
	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		for o := range output {
			fmt.Println(o)
			result = append(result, o)
		}
		outputWG.Done()
	}()

	// Close the output channel when the HTTP workers are done
	go func() {
		httpWG.Wait()
		close(output)
	}()

	// accept domains on stdin
	sc := bufio.NewScanner(os.Stdin)

	for _, d := range domains {
		fmt.Println("Scanning ", d)
		domain := strings.ToLower(d)

		// submit standard port checks
		if !skipDefault {
			httpsURLs <- domain
		}

		// Adding port templates
		xlarge := []string{"81", "300", "591", "593", "832", "981", "1010", "1311", "2082", "2087", "2095", "2096", "2480", "3000", "3128", "3333", "4243", "4567", "4711", "4712", "4993", "5000", "5104", "5108", "5800", "6543", "7000", "7396", "7474", "8000", "8001", "8008", "8014", "8042", "8069", "8080", "8081", "8088", "8090", "8091", "8118", "8123", "8172", "8222", "8243", "8280", "8281", "8333", "8443", "8500", "8834", "8880", "8888", "8983", "9000", "9043", "9060", "9080", "9090", "9091", "9200", "9443", "9800", "9981", "12443", "16080", "18091", "18092", "20720", "28017"}
		large := []string{"81", "591", "2082", "2087", "2095", "2096", "3000", "8000", "8001", "8008", "8080", "8083", "8443", "8834", "8888"}

		// submit any additional proto:port probes
		for _, p := range probes {
			switch p {
			case "xlarge":
				for _, port := range xlarge {
					httpsURLs <- fmt.Sprintf("%s:%s", domain, port)
				}
			case "large":
				for _, port := range large {
					httpsURLs <- fmt.Sprintf("%s:%s", domain, port)
				}
			default:
				pair := strings.SplitN(p, ":", 2)
				if len(pair) != 2 {
					continue
				}

				// This is a little bit funny as "https" will imply an
				// http check as well unless the --prefer-https flag is
				// set. On balance I don't think that's *such* a bad thing
				// but it is maybe a little unexpected.
				if strings.ToLower(pair[0]) == "https" {
					httpsURLs <- fmt.Sprintf("%s:%s", domain, pair[1])
				} else {
					httpURLs <- fmt.Sprintf("%s:%s", domain, pair[1])
				}
			}
		}
	}

	// once we've sent all the URLs off we can close the
	// input/httpsURLs channel. The workers will finish what they're
	// doing and then call 'Done' on the WaitGroup
	close(httpsURLs)

	// check there were no errors reading stdin (unlikely)
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}

	// Wait until the output waitgroup is done
	outputWG.Wait()

	return result
}

func isListening(client *http.Client, url, method string) bool {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := client.Do(req)
	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	if err != nil {
		return false
	}

	return true
}
