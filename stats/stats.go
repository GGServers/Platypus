package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gmemstr/platypus/common"
	"github.com/gmemstr/platypus/pluginhandler"
	"github.com/gorilla/websocket"
)

type UsageStats struct {
	Type     string  `json:"type"`
	Hostname string  `json:"hostname"`
	Cpu      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
	Disk     float64 `json:"disk"`
	Secret   string  `json:"secret"`
}

type Server struct {
	Stats  UsageStats `json:"stats"`
	Type   string     `json:"type"`
	Custom string     `json:"custom"`
	Online bool       `json:"online"`
}

var Servers map[string]Server
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handles individual connections, parsing JSON data from client and writing a message.
func Handler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		hostname := ""
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		defer c.Close()

		// @TODO Set up escalation process when socket closes.
		c.SetCloseHandler(CloseHandler)

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				_ = SetOffline(hostname)
				break
			}

			// Fire off incoming data hook.
			modifiedJson := pluginhandler.ExecuteHook(string(message), "IncomingData")
			message = []byte(modifiedJson)

			stats := UsageStats{}
			err = json.Unmarshal(message, &stats)
			if err != nil {
				break
			}
			hostname = stats.Hostname
			secretKey, err := ioutil.ReadFile(".secret")
			key := string(secretKey)
			if stats.Secret != key {
				_ = c.WriteMessage(mt, []byte("invalid secret key"))
				_ = c.Close()
			}
			// Blank out secret key after comparing.
			stats.Secret = ""
			err = WriteStats(hostname, stats)
			if err != nil {
				break
			}
			// Fire off outgoing data hook.
			reply := pluginhandler.ExecuteHook("", "OutgoingData")

			err = c.WriteMessage(mt, []byte(reply))
			if err != nil {
				break
			}
		}

		return nil
	}
}

func SetOffline(hostname string) error {
	server, ok := Servers[hostname]
	if ok {
		server.Online = false
		Servers[hostname] = server
	}

	jsonServers, err := json.MarshalIndent(Servers, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("stats.json", jsonServers, 0644)
	if err != nil {
		return err
	}
	_ = pluginhandler.ExecuteHook(hostname, "Offline")

	return nil
}

func WriteStats(hostname string, stats UsageStats) error {
	server, ok := Servers[hostname]
	if ok {
		server.Type = "stats"
		server.Stats = stats
		server.Online = true
	}
	if !ok {
		server = Server{
			Type:   "stats",
			Stats:  stats,
			Online: true,
		}
	}
	Servers[hostname] = server

	return nil
}

// Close handler, begin chain of escalation.
func CloseHandler(code int, text string) error {
	fmt.Println("ws closed")
	if code != 1000 {
		return errors.New("websocket closed badly")
	}
	return nil
}
