package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Linggle struct {
	Query  string
	Ngarms [][]interface{} `json:"ngrams"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("linggle-go [query]")
		return
	}
	query := strings.Join(os.Args[1:], " ")
	resp, err := http.Get("https://linggle.com/api/ngram/" + query)
	if err != nil {
		fmt.Println("Network failure of query: ", query)
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("IO failure of query:", query)
		panic(err)
	}
	var linggleResult Linggle
	err = json.Unmarshal(data, &linggleResult)
	if err != nil {
		fmt.Println("Json failure of query: ", query)
		panic(err)
	}
	ngrams := linggleResult.Ngarms
	if len(ngrams) == 0 {
		fmt.Println("No result found for query: ", query)
		return
	}
	totalCount := 0
	for _, item := range ngrams {
		totalCount += int(item[1].(float64))
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Results", "Counts", "Percentages"})
	for _, item := range ngrams {
		word := item[0].(string)
		count := int(item[1].(float64))
		countString := strconv.Itoa(count)
		percent := float32(count) / float32(totalCount) * 100.0
		percentString := fmt.Sprintf("%.2f", percent) + "%"
		aRecord := []string{word, countString, percentString}
		table.Append(aRecord)
	}
	table.Render()
	return
}
