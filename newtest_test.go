package newtest


import "testing"


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
	t.Log("getReport function should get the report from elasticsearch")
	{
		t.Logf("cheking getreport")
		{
			q = Query{search_string: "arpit", search_field: "user"}
			client, err := getClient()
				if err != nil {
					t.Logf("\tError getting client", ballotX)
				}
				getReport(client)
				t.Log("\tGot report successfully", checkMark)
		}
	}
}


//Benchmark getclient  function
func BenchmarkGetclient(b *testing.B) {
     b.ResetTimer()
     for i := 0; i < b.N; i++ {
         getClient()
     }
}

//Benchmark getwriter function
func BenchmarkGetcsv(b *testing.B) {
     b.ResetTimer()
     for i := 0; i < b.N; i++ {
         getwriter()
     }
}


