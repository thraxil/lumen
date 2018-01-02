package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	Port         int
	Secret       string `envconfig:"SECRET"`
	MaxProcs     int    `envconfig:"MAX_PROCS"`
	ReadTimeout  int    `envconfig:"READ_TIMEOUT"`
	WriteTimeout int    `envconfig:"WRITE_TIMEOUT"`
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

	s := Server{config: c}
	r := mux.NewRouter()

	r.HandleFunc("/", s.Index).Methods("GET")
	r.HandleFunc("/", s.Upload).Methods("POST")
	r.HandleFunc("/favicon.ico", s.Favicon)

	hs := http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		Handler:      logTop(r, sl),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		hs.ListenAndServe()
	}()

	// graceful shutdown on SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = hs.Shutdown(sctx); err != nil {
		sl.Log("level", "ERR", "msg", "error on  shutdown", "error", err)
	} else {
		sl.Log("level", "INFO", "msg", "Server Stopped")
	}
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
