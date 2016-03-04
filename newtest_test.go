package newtest


import "testing"
import "fmt"
import "gopkg.in/olivere/elastic.v3"


const checkMark =  "\u2713"
const ballotX = "\u2717"
var tr Tweet = Tweet{User:"arpit",Post_date:"2011-998-990",Message:"hello world"}
var data chan Tweet = make(chan Tweet)  
var k Fpair
var str DatabaseType

func TestCheckError(t *testing.T) {
	t.Log("Checking the error check function")
	{
		var err error
		err = nil
		CheckError(err)
		t.Logf("Error check successfull",checkMark)  
	}
}

func TestGetWriter(t *testing.T) {
	t.Log("Given the need to return a csv writer")
	{
		t.Logf("checking for function(getWriter) return")
		{
			_,err := GetWriter()
			if err != nil  {
				t.Fatal("\t Should have been a *csv.Writer",ballotX)
			}
			fmt.Println("boss")
			t.Log("\tGot the Csv Writer", checkMark)
		}
	}
}
func TestGetClient(t *testing.T) {
	t.Log("Given the need for elastic client")
	{
		t.Logf("checking return of function getClient")
		{
			str = Elasticsearch
			_,err := str.GetClient()
			if err != nil {
				t.Fatal("\tShould have been a elastic search client",ballotX)
			}
			t.Log("\tGot the elastic search client",checkMark)
		}
	}
}

func TestCsvWriter(t *testing.T) {
	t.Log("function csvWriter should  to file") 
	{
		t.Logf("checking function for return") 
		{
			k,err := GetWriter()
			if err != nil {
			t.Fatal("\t should have been a csvwriter", ballotX)
			}
			go tr.CsvWriter(k,data)
			data <- tr
			t.Log("\t data successfully written to csv", checkMark)

		}
	}
}

func TestGetField(t *testing.T) {
	t.Log("Function Getfield gets the field of struct")
	{
		t.Logf("Checking function for return")
		{
			k.Qkey = "User"
			k.Qvalue = "arpit"
			q = Filter {filter :[]Fpair{k}}
			str := GetField(&tr, q.filter[0].Qkey)
			fmt.Println(str)
			if str != "arpit" {
				t.Fatal("\t Could have got the struct field",ballotX)
			}
			t.Log("\t Field successfully got",checkMark)
		}
	}
}

func TestGetReport(t *testing.T) {
	t.Log("Function GetReport gets the report")
	{
		t.Logf("Checking function")
		{
			k.Qkey = "User"
			k.Qvalue = "arpit"
			q = Filter {filter :[]Fpair{k}}
			str = Elasticsearch
			client ,err := str.GetClient()
			if err != nil {
				t.Fatal("\tShould have been a elastic search client GetReport",ballotX)
			}
			t.Log("\tGot the elastic search client GetReport",checkMark)
			switch v := client.(type) {
				case *elastic.Client:
				fmt.Println("Calling With Elasticsearch Client")
				GetReport(v)
			default:
				fmt.Println("No such Client available",v )
			}
		}
	}
}

//Benchmark getclient  function
func BenchmarkGetclient(b *testing.B) {
     //b.ResetTimer()
	 str = Elasticsearch
     for i := 0; i < b.N; i++ {
         str.GetClient()
     }
}

//Benchmark getwriter function
func BenchmarkGetcsv(b *testing.B) {
     //b.ResetTimer()
     for i := 0; i < b.N; i++ {
         GetWriter()
     }
}

func BenchmarkGetreport(b *testing.B) {
     //b.ResetTimer()
	k.Qkey = "user"
	k.Qvalue = "arpit"
	q = Filter {filter :[]Fpair{k}}
	str = Elasticsearch
    for i := 0; i < b.N; i++ {
     	client, err := str.GetClient()
		if err != nil {
			fmt.Println("\tError getting client", ballotX)
		}
        switch v := client.(type) {
		case *elastic.Client:
			fmt.Println("Calling With Elasticsearch Client")
			GetReport(v)
		default:
			fmt.Println("No such Client available",v )
		}
     }
}

func BenchmarkCheckError(b *testing.B) {
	var err error
	err = nil 
	for i := 0; i < b.N; i++ {
		CheckError(err)
	}
}
