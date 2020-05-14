package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gmemstr/platypus/common"
	"github.com/gmemstr/platypus/pluginhandler"
	"github.com/gmemstr/platypus/stats"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handle(handlers ...common.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rc := &common.RouterContext{}
		for _, handler := range handlers {
			err := handler(rc, w, r)
			if err != nil {
				log.Printf("%v", err)

				w.Write([]byte(http.StatusText(err.StatusCode)))

				return
			}
		}
	})
}

func Init() *mux.Router {

	r := mux.NewRouter()

	// "Static" paths
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Paths that require specific handlers
	r.Handle("/", Handle(
		rootHandler(),
	)).Methods("GET")

	r.Handle("/stats", Handle(
		stats.Handler(),
	)).Methods("GET")

	r.Handle("/getstats", Handle(
		StatsWs(),
	)).Methods("GET")

	r.Handle("/poll", Handle(
		handlePoll(),
	)).Methods("GET", "POST")

	return r
}

func StatsWs() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		s := stats.Servers
		err = c.WriteJSON(s)
		if err != nil {
			c.Close()
		}
		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					break
				}
				log.Printf("recv: %v", message)
			}
		}()

		interval := time.Duration(common.Config.Interval) * time.Second
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return nil
			case <-ticker.C:
				s := stats.Servers
				err = c.WriteJSON(s)
				if err != nil {
					break
				}
			}
		}

		return nil
	}
}

func rootHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		var file string
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			file = "web/index.html"
		default:
			return &common.HTTPError{
				Message:    fmt.Sprintf("%s: Not Found", r.URL.Path),
				StatusCode: http.StatusNotFound,
			}
		}

		return common.ReadAndServeFile(file, w)
	}
}

func handlePoll() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		authkey := r.Header.Get("X-PLATYPUS-AUTH")
		appType := r.Header.Get("X-PLATYPUS-TYPE")

		secretKey, err := ioutil.ReadFile(".secret")
		if err != nil {
			return nil
		}
		key := string(secretKey)
		if key != authkey {
			return nil
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		payload := buf.String()

		processed := pluginhandler.ExecuteHook(payload, "IncomingData")
		s := stats.Server{Custom: processed, Type: appType}

		stats.Servers[appType] = s

		response := pluginhandler.ExecuteHook("", "OutgoingData")
		fmt.Fprintln(w, response)
		jsonServers, err := json.MarshalIndent(stats.Servers, "", "  ")
		if err != nil {
			return nil
		}
		err = ioutil.WriteFile("stats.json", jsonServers, 0644)
		if err != nil {
			return nil
		}

		return nil
	}
}

