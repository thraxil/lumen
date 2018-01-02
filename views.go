package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	config config
}

func (s Server) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
<head>
<title>Lumen: Upload Image</title>
<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css" />
</head>

<body>
<div class="container">
<h1>Upload Image</h1>

<form action="." method="post" enctype="multipart/form-data" >

<input type="file" name="image" class="form-control"/><br />
secret: <input type="password" name="secret" class="form-control"/><br />
<input type="submit" value="upload image" class="btn btn-primary"/>
</form>
</div>
</body>
</html>`)
}

func (s Server) Upload(w http.ResponseWriter, r *http.Request) {

}

func (s Server) Favicon(w http.ResponseWriter, r *http.Request) {
	// just ignore this crap
}
