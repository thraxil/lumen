package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/thraxil/resize"
)

type Server struct {
	config  config
	backend backend
	logger  log.Logger
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

var mimeexts = map[string]string{
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
	"image/png":  ".png",
}

var extmimes = map[string]string{
	".jpg": "image/jpeg",
	".gif": "image/gif",
	".png": "image/png",
}

type imageData struct {
	Hash      string `json:"hash"`
	Length    int64  `json:"length"`
	Extension string `json:"extension"`
	FullURL   string `json:"full_url"`
	// we don't use these fields, but retain them for reticulum compat
	Satisfied bool     `json:"satisfied"`
	Nodes     []string `json:"nodes"`
}

func (s Server) Upload(w http.ResponseWriter, r *http.Request) {
	if !s.config.ValidKey(r.FormValue("secret")) {
		// TODO:  log this
		http.Error(w, "invalid upload secret", 403)
		return
	}

	i, fh, _ := r.FormFile("image")
	defer i.Close()
	h := sha1.New()
	io.Copy(h, i)
	ahash, err := hashFromString(fmt.Sprintf("%x", h.Sum(nil)), "")
	if err != nil {
		http.Error(w, "bad hash", 500)
		return
	}
	i.Seek(0, 0)

	mimetype := fh.Header.Get("Content-Type")
	if mimetype == "" {
		// they left off a mimetype, so default to jpg
		mimetype = "image/jpeg"
	}
	ext := mimeexts[mimetype]

	ri := imageSpecifier{
		ahash,
		resize.MakeSizeSpec("full"),
		ext,
	}
	err = s.backend.Write(ri.fullSizePath(), i)
	if err != nil {
		http.Error(w, "S3: "+err.Error(), 500)
		return
	}

	id := imageData{
		Hash:      ahash.String(),
		Extension: ext,
		FullURL:   "/image/" + ahash.String() + "/full/image" + ext,
		Satisfied: true,
		Length:    fh.Size,
		Nodes:     []string{},
	}
	b, err := json.Marshal(id)
	if err != nil {
		s.logger.Log("level", "ERR", "error", err.Error())
	}
	w.Write(b)
	//  update metrics/dashboard
	//	s.Uploaded(imageRecord{*ahash, ext})
}

func (s Server) Serve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	size := vars["size"]
	filename := vars["filename"]

	ahash, err := hashFromString(hash, "")
	if err != nil {
		http.Error(w, "invalid hash", http.StatusNotFound)
		return
	}

	sspec := resize.MakeSizeSpec(size)
	if sspec.String() != size {
		// force normalization of size spec
		http.Redirect(w, r, "/image/"+ahash.String()+"/"+sspec.String()+"/"+filename, http.StatusMovedPermanently)
		return
	}

	extension := filepath.Ext(filename)
	if extension == ".jpeg" {
		// normalize .jpeg to .jpg
		fixedFilename := strings.Replace(filename, ".jpeg", ".jpg", 1)
		http.Redirect(w, r, "/image/"+ahash.String()+"/"+sspec.String()+"/"+fixedFilename, http.StatusMovedPermanently)
		return
	}
	ri := &imageSpecifier{ahash, sspec, extension}

	imgData, err := s.backend.Read(ri.Path())
	if err != nil {
		// for now we just have to 404
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else {
		w = setCacheHeaders(w, ri.Extension)
		w.Write(imgData)
	}

	fmt.Println(hash, size, filename)
}

func (s Server) Favicon(w http.ResponseWriter, r *http.Request) {
	// just ignore this crap
}

func setCacheHeaders(w http.ResponseWriter, extension string) http.ResponseWriter {
	w.Header().Set("Content-Type", extmimes[extension])
	w.Header().Set("Expires", time.Now().Add(time.Hour*24*365).Format(time.RFC1123))
	return w
}
