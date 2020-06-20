package model

import (
	"bufio"
	"strings"

	// "io/ioutil"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/nachol/scanner/utils"
)

var DirsearchScanresult []string
var DirsearchScanraw string

func DirsearchScan(s *Scan, args ...interface{}) ([]string, string, error) {

	for _, t := range s.Scope {
		log.Println(t)
		output := "/tmp/dirsearch.txt"

		cmd := exec.Command("python3", "./tools/dirsearch/dirsearch.py", "-u", t, "-E", "--json-report", output, "-t", strconv.Itoa(s.Threads))
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
		err = DirsearchformatResult(output)
		if err != nil {
			log.Println(err.Error())
			return nil, "", err
		}
		DirsearchScanraw = DirsearchScanraw + out.String()
	}
	log.Println(DirsearchScanraw)
	return utils.Unique(DirsearchScanresult), DirsearchScanraw, nil

}

func DirsearchformatResult(output string) error {
	file, err := os.Open(output)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		DirsearchScanresult = append(SubdomainScanresult, strings.TrimSpace(scanner.Text()))
	}
	return nil
}
