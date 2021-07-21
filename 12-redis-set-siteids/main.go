package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"redissetsiteids/connections"
	"strings"
)

func main() {
	// Read file from command line arguments
	siteids := flag.String("siteids", "", "File containing siteids")
	flag.Parse()

	file, err := os.Open(*siteids)
	if err != nil {
		fmt.Println("Some issue one file reading.")
		log.Fatal(err)
	}
	defer file.Close()
	r := bufio.NewScanner(file)

	client, err := connections.RedisConnect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	for r.Scan() {
		siteid := r.Text()
		splitSiteID := strings.Split(siteid, " ")
		key := strings.Trim(splitSiteID[1], "\"")
		value := strings.Trim(splitSiteID[2], "\"")
		fmt.Println(key, value)
		// fmt.Println(client.Set(key, value, 0).Val())
	}
	fmt.Println(os.Getenv("USER"))
}
