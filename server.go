package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

var Version string

type Server struct {
	clientAddress string
}

// Apilisten has to be made dynamic for launching different url for different server
func (s *Server) Apilisten(response *Response) {
	address := s.clientAddress

	Version = checkGitCommitVersion()
	mux := http.NewServeMux()
	mux.HandleFunc("/hosts", response.GetTargets)
	mux.HandleFunc("/healthcheck", HealthCheck)
	mux.HandleFunc("/version", CheckVersion)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	err := http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Debug("Sending HealthCheck Success")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func CheckVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Version ", Version)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(Version))
}

func checkGitCommitVersion() string {
	var Version = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}

		return ""
	}()
	return Version
}

func (re *Response) GetTargets(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(*re)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(200)
	w.Write(res)
}
