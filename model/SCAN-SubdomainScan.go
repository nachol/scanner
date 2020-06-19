package model

import (
	"bufio"
	// "io/ioutil"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/nachol/scanner/utils"
)

var SubdomainScanresult []string
var SubdomainScanraw string

func SubdomainScan(s *Scan, args ...interface{}) ([]string, string, error) {

	for _, domain := range s.Scope {
		log.Println(domain)
		output := "/tmp/sublist3r.txt"

		cmd := exec.Command("python3", "./tools/Sublist3r/sublist3r.py", "-d", domain, "-o", output, "-t", strconv.Itoa(s.Threads))
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Start()
		if err != nil {
			log.Println(err.Error())
			return nil, "", err
		}
		err = cmd.Wait()
		if err != nil {
			log.Println(err.Error())
			return nil, "", err
		}
		err = SubdomainformatResult(output)
		if err != nil {
			log.Println(err.Error())
			return nil, "", err
		}
		SubdomainScanraw = SubdomainScanraw + "\n" + out.String()
	}
	log.Println(SubdomainScanraw)
	return utils.Unique(SubdomainScanresult), SubdomainScanraw, nil

}

func SubdomainformatResult(output string) error {
	file, err := os.Open(output)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		SubdomainScanresult = append(SubdomainScanresult, scanner.Text())
	}
	return nil
}

//
// Returns unique items in a slice
//
func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}
