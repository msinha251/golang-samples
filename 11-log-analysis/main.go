package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

func main() {

	// read arguments from the command line
	path := flag.String("path", "sample.log", "path to the log file")
	level := flag.String("level", "ERROR", "log level")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, *level) {
			log.Println(line)
		}
	}
}
