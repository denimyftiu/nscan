package main

import (
	"bufio"
	"flag"
	"fmt"
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

func scanDomain(domain string, options string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Create("./nmap/" + domain)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	args := strings.Split(fmt.Sprintf("nmap %s %s", options, domain), " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = f

	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	var path string
	var options string
	flag.StringVar(&path, "path", "domains", "-path <path to domains file>")
	flag.StringVar(&options, "options", "", "-options <nmap scan options>")
	flag.Parse()

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist.")
			return
		}
	}

	dir := filepath.Dir(path)
	var err error
	err = os.Chdir(dir)
	err = os.Mkdir("nmap", 0755)
	if err != nil {
		fmt.Println("Directory exists.")
		return
	}

	domains, err := getDomains(path)
	if err != nil {
		fmt.Println("Can not read file.")
		return
	}
	var wg sync.WaitGroup
	for _, d := range domains {
		wg.Add(1)
		go scanDomain(d, options, &wg)
	}
	wg.Wait()
}
