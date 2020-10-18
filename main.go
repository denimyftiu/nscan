package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func getDomains(filePath string) ([]string, error) {

	var domains []string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return domains, nil
}

func scanDomain(domain string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Create("./nmap/" + domain)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := "nmap -vv -sV -sC " + domain
	args := strings.Split(s, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = f

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	var path string
	flag.StringVar(&path, "path", "domains", "-path <path to domains file>")
	flag.Parse()

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			panic(err)
		}
	}

	dir := filepath.Dir(path)
	var err error
	err = os.Chdir(dir)
	err = os.Mkdir("nmap", 0755)
	if err != nil {
		panic(err)
	}

	domains, err := getDomains(path)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, d := range domains {
		wg.Add(1)
		go scanDomain(d, &wg)
	}
	wg.Wait()
}
