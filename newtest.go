package main 

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"os"
	"sync"
	"reflect"
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

type Filter struct {
	filter []Fpair
}

type Fpair struct {
	Qkey  string
	Qvalue string
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
		//fmt.Println(c)
		data := []string{c.User, c.Post_date, c.Message}
		//Introduced locks for write to csv file
		mutex.Lock() 
		writer.Write(data)
		writer.Flush()
		mutex.Unlock()
		//lock closed
	}
}

//Now Generic lookup is possible thanks to this function
func getField(v *Tweet, field string) string {
    r := reflect.ValueOf(v)
    f := reflect.Indirect(r).FieldByName(field)
   // fmt.Println(string(f.String()))
    return string(f.String())
}

func filtering(search chan *elastic.SearchResult)  {
	var t Tweet
	var data chan Tweet = make(chan Tweet)
	//var filter string 
	//fmt.Println("csv writer started")
	writer, err := getwriter()
	if err != nil {
		panic(err)
	}
	go t.csvWriter(writer, data)  // spawning the csvwriter routine
	//fmt.Println("filtering started")
	for i := range search {
		searchResult := i
		for _, hit := range searchResult.Hits.Hits {
	        err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				fmt.Println("failed", err)
			}
			//Filtering goes her 
			//val := reflect.ValueOf(q.filter[0].Qkey)

			//fmt.Println(getField(&t, q.filter[0].Qkey))
			q.filter[0].Qkey = strings.Replace(q.filter[0].Qkey,q.filter[0].Qkey[:1], strings.ToUpper(q.filter[0].Qkey[:1]),1)
			if getField(&t, q.filter[0].Qkey) == q.filter[0].Qvalue {
				fmt.Println(t)
				data <- t
			}
		}
	}
	close(data) // closing the channel
}

func getReport(client *elastic.Client) {
	result := make(chan *elastic.SearchResult)
	// spawinng the filtering routine
	go filtering(result) 
	// the termquery uses all lower but for matching to filter exactly we have to convert the first letter to upper
	boolq := elastic.NewBoolQuery()
	termQuery := boolq.Filter(elastic.NewTermQuery(q.filter[0].Qkey, q.filter[0].Qvalue))
	count, err := client.Count().
		Query(termQuery).
		Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Count", count)
	scrollService := elastic.NewScrollService(client)
	searchResult, err := scrollService.Scroll("5m").Size(1).Do()
	if err != nil {
		panic(err)
	}
	pages := 0
	scroll_indexId := searchResult.ScrollId
	for {
		searchResult, err := scrollService.Query(termQuery).Scroll("5m").
			Size(1).
			ScrollId(scroll_indexId).
			Do()
		if err != nil {
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
	close(result)  //closing the channel

}

var q Filter // Global because it has to be used at different routines
func main() {
	var k Fpair
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Enter the search Field")
	fmt.Scan(&k.Qkey)

	fmt.Println("Enter the search string")
	fmt.Scan(&k.Qvalue)

	q = Filter {filter :[]Fpair{k}}
	fmt.Println(q.filter[0])
	getReport(client)
}
