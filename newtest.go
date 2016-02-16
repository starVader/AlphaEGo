package main 

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"os"
	"sync"
	"strings"
)

type Tweet struct {
	User      string
	Post_date string
	Message   string
}

type output interface {
	csvWriter()
}

type Query struct {
	search_string string
	search_field  string
}

func getClient() (*elastic.Client, error) {
	client, err := elastic.NewClient()
	return client, err
}

func getwriter() (*csv.Writer, error) {
	file, err := os.Create("result.csv")
	writer := csv.NewWriter(file)
	return writer, err
}

func (c Tweet) csvWriter(writer *csv.Writer, m chan Tweet) {
	var mutex = &sync.Mutex{}
	for i := range m {
		c = i
		fmt.Println(c)
		data := []string{c.User, c.Post_date, c.Message}
		mutex.Lock() //Introduced locks for write to csv file
		writer.Write(data)
		writer.Flush()
		mutex.Unlock()
	}
}

func filtering(search chan *elastic.SearchResult) {
	var t Tweet
	var data chan Tweet = make(chan Tweet)
	fmt.Println("csv writer started")
	writer, err := getwriter()
	if err != nil {
		panic(err)
	}
	go t.csvWriter(writer, data)  // spawning the csvwriter routine
	fmt.Println("filtering started")
	for i := range search {
		searchResult := i
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				fmt.Println("failed", err)
			}
			filter := strings.Replace(q.search_string,q.search_string[:1], strings.ToUpper(q.search_string[:1]),1) // filtering according to searchstring
			if t.User == filter {
				data <- t // channeling data to csv writer
			}
		}
	}
	close(data)
}

func getReport(client *elastic.Client) {
	result := make(chan *elastic.SearchResult)
	go filtering(result) // spawinng the filtering routine
	termQuery := elastic.NewTermQuery(q.search_field, q.search_string)// the termquery uses all lower but for matching to filter exactly we have to convert the first letter to upper
	count, err := client.Count().
		Query(termQuery).
		Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Count", count)
	searchResult, err := client.Scroll().Size(1).Do()
	if err != nil {
		panic(err)
	}
	pages := 0
	scroll_indexId := searchResult.ScrollId
	data := make(chan Tweet)
	for {
		searchResult, err := client.Scroll().
			Size(1).
			ScrollId(scroll_indexId).
			Do()
		if err != nil {
			fmt.Println(err)
			break
		}
		result <- searchResult // sending data into channel received by filtering function
		pages += 1
		scroll_indexId = searchResult.ScrollId
		if scroll_indexId == "" {
			fmt.Println(scroll_indexId)
		}
	}

	if pages <= 0 {
		fmt.Println(pages, "Records found")
	}
	close(data)

}

var q Query // Global because it has to be used at different routines
func main() {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Successfully got the client:", client)
	fmt.Println("Enter the search Field")
	fmt.Scan(&q.search_field)
	fmt.Println("Enter the search string")
	fmt.Scan(&q.search_string)
	getReport(client)
}
