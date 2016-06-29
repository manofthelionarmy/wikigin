package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var responseTemplate = template.Must(template.New("results").Parse(`
	<html>
	<head>
		<title>results</title>
	</head>
	<body>
		<ol>
		{{range .Results}}
			<li>{{.Title}} -<a href="{{.URL}}">{{.URL}}</a></li>
			{{end}}
		</ol>
		<p>{{len .Results}} results in {{.Elapsed}}</p>
	</body>
	</html>

 `))

//Creating a response struct
//Maybe there be packages that'll have their own Result types
type Result struct {
	Title, URL string
}

//structure the search results
type Response struct {
	Results []Result

	Elapsed time.Duration
}

func main() {

	/*http.HandleFunc("/hello", handleHello)

	fmt.Print("serving on http://localhost:7777/hello")

	log.Fatal(http.ListenAndServe("localhost:7777", nil))*/

	http.HandleFunc("/search", handleSearch)

	fmt.Print("serving on http://localhost:7777/search")

	log.Fatal(http.ListenAndServe("localhost:7777", nil))
}

func handleHello(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)
	fmt.Fprintln(w, "<h1>hello world</h1>")
}

//validate the search query
func handleSearch(w http.ResponseWriter, req *http.Request) {

	var err error
	log.Println("serving", req.URL)

	query := req.FormValue("q")

	if query == "" {
		http.Error(w, `missing "q" URL parameter`, http.StatusBadRequest)
	}

	//fetch the search results

	start := time.Now()

	//eventually need to check if there is an error with the Search(query)

	elapsed := time.Since(start)

	results := make([]Result, 3)

	results[0].Title = "YouTube"
	results[0].URL = "http://www.youtube.com"

	results[1].Title = "Google"

	results[1].URL = "http://www.google.com"

	results[2].Title = "Facebook"

	results[2].URL = "http://www.facebook.com"

	//Render the search results

	resp := Response{results, elapsed}

	switch req.FormValue("output") {
	case "json":
		err = json.NewEncoder(w).Encode(resp)
	case "prettyjson":
		var b []byte //instead of storing it by string, whcih will take up to much memory, it is better to store its byte value
		b, err = json.MarshalIndent(resp, "", "  ")
		if err == nil {
			_, err = w.Write(b)
		}
	default: //HTML
		err = responseTemplate.Execute(w, resp)
	}

}
