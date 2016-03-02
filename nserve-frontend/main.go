package main

import (
	"encoding/json"
	. "github.com/THUNDERGROOVE/nserve/lib"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var tmplText = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
<title>Maximumtwang status</title>

<!-- Bootstrap -->
<!-- Latest compiled and minified CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">

<!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
<!--[if lt IE 9]>
																						<script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
																									<script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
																											<![endif]-->
</head>
<body>

<div class="navbar navbar-inverse navbar-static-top">
	<div class="navbar-header">
		<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1" aria-expanded="false">
	  <span class="sr-only">Toggle navigation</span>
	  <span class="icon-bar"></span>
	  <span class="icon-bar"></span>
	  <span class="icon-bar"></span>
	  </button>
	  <a class="navbar-brand" href="maximumtwang.com">Maximumtwang</a>
  </div>
	<ul class="nav navbar-nav">
		<li class="active"><a href="status.maximumtwang.com">status</a></li>
		<li><a href="eve.maximumtwang.com">EVEMuGate</a></li>

	</ul>
</div>

<div class="container">
	{{if .HasError}}
		<h1>Error!!!! :(</h1>
		<p>{{.Error.Error}}</p>
	{{else}}
		<table class="table table-striped">
		<thead>
			<tr>
				<th>Service</th>
				<th>Status</th>
			</tr>
		</thead>
		<tbody>
	{{range .Targets}}
		<tr>
			<td>{{.Name}}</td>
			{{if .Running}}
				<td><span class="label label-success">Okay</span></td>
			{{else}}
				<td><span class="label label-danger">Down :(</span></td>
			{{end}}
		</tr>
	{{end}}
		</tbody>
	</table>
	{{end}}
	</div>

	<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
	<!-- Include all compiled plugins (below), or include individual files as needed -->
	</body>
	</html>
	`

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

	//mainT = template.Must(template.ParseFiles("main.tmpl"))
	mainT = template.New("")
	mainT = template.Must(mainT.Parse(tmplText))

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
