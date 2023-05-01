package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/EasyRecon/wappaGo/cmd"
	"github.com/EasyRecon/wappaGo/structure"
	"github.com/EasyRecon/wappaGo/technologies"
)

func main() {
	options := structure.Options{}
	options.Screenshot = flag.String("screenshot", "", "path to screenshot if empty no screenshot")
	options.Ports = flag.String("ports", "80,443", "port want to scan separated by coma")
	options.Threads = flag.Int("threads", 5, "Number of threads to start recon in same time")
	options.Report = flag.Bool("report", false, "Generate HTML report")
	options.Porttimeout = flag.Int("port-timeout", 2000, "Timeout during port scanning in ms")
	//options.ChromeTimeout = flag.Int("chrome-timeout", 0000, "Timeout during navigation (chrome) in sec")
	options.ChromeThreads = flag.Int("chrome-threads", 5, "Number of chromes threads in each main threads total = option.threads*option.chrome-threads (Default 5)")
	options.Resolvers = flag.String("resolvers", "", "Use specifique resolver separated by comma")
	options.AmassInput = flag.Bool("amass-input", false, "Pip directly on Amass (Amass json output) like amass -d domain.tld | wappaGo")
	options.FollowRedirect = flag.Bool("follow-redirect", false, "Follow redirect to detect technologie")
	options.Proxy = flag.String("proxy", "", "Use http proxy")
	flag.Parse()
	configure(options)
}

func configure(options structure.Options) {
	if *options.Screenshot != "" {
		if _, err := os.Stat(*options.Screenshot); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(*options.Screenshot, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
	}
	folder, errDownload := technologies.DownloadTechnologies()
	if errDownload != nil {
		log.Println("error during downloading techno file")
	}
	defer os.RemoveAll(folder)

	var input []string
	var scanner = bufio.NewScanner(bufio.NewReader(os.Stdin))
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	c := cmd.Cmd{}
	c.ResultGlobal = technologies.LoadTechnologiesFiles(folder)
	c.Options = options
	c.Input = input

	results := make(chan structure.Data)

	go func() {
		for result := range results {
			b, err := json.Marshal(result)

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(b))
		}
	}()

	c.Start(results)
}
