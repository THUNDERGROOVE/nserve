package main

import (
	"encoding/json"
	. "github.com/THUNDERGROOVE/nserve/lib"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var mainT *template.Template

type Context struct {
	Targets  []Target
	HasError bool
	Error    error
}

func CheckError(err error, c *Context, rw http.ResponseWriter) bool {
	if err != nil {
		c.HasError = true
		c.Error = err

		err := mainT.Execute(rw, *c)
		if err != nil {
			log.Printf("Failed to execute template: %s\n", err.Error())
		}
		return true
	}

	return false
}

func main() {

	mainT = template.Must(template.ParseFiles("main.tmpl"))

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		var c Context
		var targets []Target
		resp, err := http.Get("http://localhost:5598/")
		if CheckError(err, &c, rw) {
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if CheckError(err, &c, rw) {
			return
		}

		err = json.Unmarshal(data, &targets)

		c.Targets = targets

		if CheckError(err, &c, rw) {
			return
		}
		err = mainT.Execute(rw, c)
		if err != nil {
			log.Printf("Failed to execute template: %s\n", err.Error())
		}
	})
	http.ListenAndServe(":5599", nil)
}
