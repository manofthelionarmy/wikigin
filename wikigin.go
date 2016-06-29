package main

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	//fmt.Println(title + "\n" + body)
	return &Page{Title: title, Body: body}, nil
}

func handler(c *gin.Context) {
	p := c.Param("title")
	fmt.Fprintf(c.Writer, "<h1>Hi %s</h1>", p)
}

func viewHandler(c *gin.Context) {
	title := c.Param("title")
	p, _ := loadPage(title)

	fmt.Fprintf(c.Writer, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(c *gin.Context) {
	title := c.Param("page")
	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	t, _ := template.ParseFiles("edit.html")
	t.Execute(c.Writer, p)
}

func main() {

	router := gin.Default()

	//http.HandleFunc("/view/", viewHandler)

	//router.GET("/:title", handler)

	//http.ListenAndServe(":8080", nil)
	router.GET("/view/:title", viewHandler) //side note, cannot have same wildcard names
	router.GET("/edit/:page", editHandler)
	//router.GET("/save/:saved", saveHandler)
	router.Run(":8080")
}
