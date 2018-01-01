package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Backend string

	S3AccessKey string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey string `envconfig:"S3_SECRET_KEY"`
	S3Bucket    string `envconfig:"S3_BUCKET"`

	Port          int
	ClusterSecret string `envconfig:"CLUSTER_SECRET"`
	MaxProcs      int    `envconfig:"MAX_PROCS"`
	ReadTimeout   int    `envconfig:"READ_TIMEOUT"`
	WriteTimeout  int    `envconfig:"WRITE_TIMEOUT"`
}

func main() {
	sl := newSTDLogger()
	sl.Log("level", "INFO", "msg", "starting logger")

	var c config
	err := envconfig.Process("lumen", &c)
	if err != nil {
		sl.Log("level", "ERR", "error", err)
		os.Exit(1)
	}

	mux.NewRouter()

}

func logTop(handler http.Handler, sl log.Logger) http.Handler {
	sl = log.With(sl, "component", "web")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rc := recover(); rc != nil {
				sl.Log("level", "ERR", "path", r.URL.String(), "msg", rc)
			}
		}()
		t0 := time.Now()
		handler.ServeHTTP(w, r)
		t1 := time.Now()
		sl.Log(
			"level", "INFO",
			"remote_addr", r.RemoteAddr,
			"method", r.Method,
			"path", r.URL.String(),
			"time",
			fmt.Sprintf("%v", t1.Sub(t0)))
	})
}
