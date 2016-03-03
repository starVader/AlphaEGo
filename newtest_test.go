package newtest


import "testing"
import "fmt"
import "gopkg.in/olivere/elastic.v3"


const checkMark =  "\u2713"
const ballotX = "\u2717"
var tr Tweet = Tweet{User:"rakesh",Post_date:"2011-998-990",Message:"hello world"}
var data chan Tweet = make(chan Tweet)  


func TestFunction(t *testing.T) {
	t.Log("Given the need to return a csv writer")
	{
		t.Logf("checking for function(getWriter) return")
		{
			_,err := getwriter()
			if err != nil  {
				t.Fatal("\t Should have been a *csv.Writer",ballotX)
			}
			t.Log("\tGot the Csv Writer", checkMark)
		}
	}
	t.Log("Given the need for elastic client")
	{
		t.Logf("checking return of function getClient")
		{
			_,err := getClient()
			if err != nil {
				t.Fatal("\tShould have been a elastic search client",ballotX)
			}
			t.Log("\tGot the elastic search client",checkMark)
		}
	}
	t.Log("function csvWriter should  to file") 
	{
		t.Logf("checking function for return") 
		{
			k,err := getwriter()
			if err != nil {
				t.Fatal("\t should have been a csvwriter", ballotX)
			}
			go tr.csvWriter(k,data)
			data <- tr
			t.Log("\t data successfully written to csv", checkMark)

		}
	}
	t.Log("testing function getReport and filtering")
	{
		t.Logf("cheking getreport")
		{
			client, err := getClient()
			if err != nil {
				t.Logf("\tError getting client", ballotX)
			}
			q = Query{search_string: "arpit", search_field: "user"}
			getReport(client)
			t.Log("\tGot report successfully", checkMark)
			k := make(chan *elastic.SearchResult)
			go filtering(k)
			searchResult, err := client.Scroll().Size(1).Do()
			if err != nil {
				t.Log(err)
			}
			scroll_indexId := searchResult.ScrollId
			for {
				searchResult, err := client.Scroll().
				Size(1).
				ScrollId(scroll_indexId).
				Do()
				if err != nil {
					break
				}
				k <- searchResult
			}
			
		}
	}
}


//Benchmark getclient  function
func BenchmarkGetclient(b *testing.B) {
     //b.ResetTimer()
     for i := 0; i < b.N; i++ {
         getClient()
     }
}

//Benchmark getwriter function
func BenchmarkGetcsv(b *testing.B) {
     //b.ResetTimer()
     for i := 0; i < b.N; i++ {
         getwriter()
     }
}

func BenchmarkGetreport(b *testing.B) {
     //b.ResetTimer()
     for i := 0; i < b.N; i++ {
     	client, err := getClient()
		if err != nil {
			fmt.Println("\tError getting client", ballotX)
		}
        getReport(client)
     }
}

