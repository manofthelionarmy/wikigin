package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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

func viewHandler(c *gin.Context, title string) {

	p, err := loadPage(title)
	if err != nil {
		c.Redirect(http.StatusFound, "/edit/"+title)
		return
	}
	renderTemplate(c, "view", p)
}

func editHandler(c *gin.Context, title string) {

	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(c, "edit", p)
}

func saveHandler(c *gin.Context, title string) {

	body := c.Request.FormValue("body")

	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		//http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/view/"+title)
}

func renderTemplate(c *gin.Context, tmpl string, p *Page) {

	err := templates.ExecuteTemplate(c.Writer, tmpl+".html", p)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func getTitle(c *gin.Context) (string, error) {
	m := validPath.FindStringSubmatch(c.Request.URL.Path)

	if m == nil {
		c.AbortWithError(http.StatusNotFound, errors.New("invalid page"))
		//http.NotFound(c.Writer, c.Request)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func makeHandler(fn func(*gin.Context, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		m := validPath.FindStringSubmatch(c.Request.URL.Path)
		if m == nil {
			c.AbortWithError(http.StatusNotFound, errors.New("Invalid Page"))
			return
		}
		fn(c, m[2])
	}
}
func main() {

	router := gin.Default()

	//http.HandleFunc("/view/", viewHandler)

	//router.GET("/:title", handler)

	//http.ListenAndServe(":8080", nil)
	router.GET("/view/:title", makeHandler(viewHandler)) //side note, cannot have same wildcard names
	router.GET("/edit/:page", makeHandler(editHandler))

	router.POST("/save/:saved", makeHandler(saveHandler))
	router.Run(":8080")
}
