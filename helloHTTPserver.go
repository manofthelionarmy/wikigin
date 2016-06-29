package main

import (
	"encoding/json"
	"fmt"
	//"golang.org/x/talks/2016/applicative/google"
	"html/template"
	"log"
	"net/http"
	"time"
)

var responseTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
  <ol>
  {{range .Results}}
    <li>{{.Title}} - <a href="{{.URL}}">{{.URL}}</a></li>
  {{end}}
  </ol>
  <p>{{len .Results}} results in {{.Elapsed}}</p>
</body>
</html>
`))

func main() {
	//http.HandleFunc("/hello", handleHello)
	//fmt.Println("serving on http://localhost:4000/hello")
	//log.Fatal(http.ListenAndServe("localhost:4000", nil))

	http.HandleFunc("/search", handleSearch)
	fmt.Println("serving on http://localhost:8080/search")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleHello(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	fmt.Fprintln(w, "<h1>hello, world!</h1>")
}

type Result struct {
	Title string
	URL   string
}

func handleSearch(w http.ResponseWriter, req *http.Request) {
	var err error

	log.Println("serving", req.URL)

	//fmt.Fprintln(w, "<h1>it's suppose to search</h1>")
	query := req.FormValue("q")

	if query == "" {
		http.Error(w, `missing "q URL parameter`, http.StatusBadRequest)

		return
	}

	start := time.Now()

	//results, err := google.Search(query)

	elapsed := time.Since(start)

	/*if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/

	r := make([]Result, 2)

	r[0].URL = "http://www.google.com"
	r[0].Title = "Google"

	r[1].URL = "http://www.facebook.com"

	r[1].Title = "Facebook"

	type response struct {
		Results []Result
		Elapsed time.Duration
	}

	resp := response{r, elapsed}

	switch req.FormValue("output") {
	case "json":
		err = json.NewEncoder(w).Encode(resp)
	case "prettyjson":
		var b []byte
		b, err = json.MarshalIndent(resp, "", " ")
		if err == nil {
			_, err = w.Write(b)
		}
	default:
		err = responseTemplate.Execute(w, resp)
	}

}
