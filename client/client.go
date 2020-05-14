package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"os"
	"time"
)

type UsageStats struct {
	Type     string  `json:"type"`
	Hostname string  `json:"hostname"`
	Cpu      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
	Disk     float64 `json:"disk"`
	Secret   string  `json:"secret"`
}

type Configuration struct {
	Master string `yaml:"master"`
	Secret string `yaml:"secret"`
}

func main() {
	conf := Configuration{}
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	addr := conf.Master
	if addr == "" {
		fmt.Println("Missing master server URL, please edit your config.yml")
		os.Exit(1)
	}
	if conf.Secret == "" {
		fmt.Println("Missing master server secret key, please edit your config.yml")
		os.Exit(1)
	}
	u := url.URL{Scheme: "ws", Host: addr, Path: "/stats"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
		}
		c.Close()
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			stats, err := GetStats()
			if err != nil {
				break
			}
			stats.Secret = conf.Secret
			statsJson, err := json.Marshal(stats)
			if err != nil {
				break
			}
			err = c.WriteMessage(websocket.TextMessage, statsJson)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

// Fetch usage stats using gopsutil.
func GetStats() (UsageStats, error) {
	stats := UsageStats{
		Type: "server",
	}
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return stats, err
	}
	stats.Disk = math.Round(diskUsage.UsedPercent)

	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return stats, err
	}
	stats.Cpu = math.Round(cpuUsage[0])

	memUsage, err := mem.VirtualMemory()
	if err != nil {
		return stats, err
	}
	stats.Memory = math.Round(memUsage.UsedPercent)

	hostname, err := os.Hostname()
	if err != nil {
		return stats, err
	}
	stats.Hostname = hostname
	return stats, nil
}
