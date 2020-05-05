package scan

import (
	"bufio"
	// "io/ioutil"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var SubdomainScanresult []string
var SubdomainScanraw string

func SubdomainScan(s *Scan, args ...interface{}) (interface{}, string, error) {

	//./sublist3r.py -d koho.ca -o ~/Documents/koho.ca/sublister
	for _, domain := range s.Scope {
		log.Println(domain)
		output := "/tmp/sublist3r.txt"

		cmd := exec.Command("python", "./tools/Sublist3r/sublist3r.py", "-d", domain, "-o", output, "-t", strconv.Itoa(s.Threads))
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
	return map[string]interface{}{
		"result": SubdomainScanresult,
	}, SubdomainScanraw, nil

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
