package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

type Fact struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

func main() {
	allFacts := make([]Fact, 0)

	collector := colly.NewCollector(colly.AllowedDomains())
	collector.OnHTML(".factsList li", func(h *colly.HTMLElement) {
		factID, err := strconv.Atoi(h.Attr("id"))
		if err != nil {
			fmt.Println("couldn't get the ID")
			log.Fatal(err)
		}
		factDesc := h.Text
		var fact Fact
		fact.Id = factID
		fact.Description = factDesc
		allFacts = append(allFacts, fact)
	})
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL.String())
	})

	collector.Visit("https://www.factretriever.com/top-10-smartest-animals")

	// // Display JSON
	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", " ")
	// enc.Encode(allFacts)

	writeJSON(allFacts)
}

func writeJSON(data []Fact) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Unable to write JSON file.")
		log.Fatal(err)
	}

	ioutil.WriteFile("Top10Facts.json", file, 0644)
}
