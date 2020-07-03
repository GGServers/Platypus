package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gmemstr/platypus/common"
	"github.com/gmemstr/platypus/pluginhandler"
	"github.com/gmemstr/platypus/router"
	"github.com/gmemstr/platypus/stats"
	"github.com/go-yaml/yaml"
)

func main() {
	genFiles()
	pluginhandler.RegisterPlugins()

	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &common.Config)
	if err != nil {
		panic(err)
	}

	// Repopulate servers from cache file.
	statsCache, err := ioutil.ReadFile("stats.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(statsCache, &stats.Servers)
	if err != nil {
		panic(err)
	}
	// Cache server stats in event master goes down.
	defer func() {
		jsonServers, err := json.MarshalIndent(stats.Servers, "", "  ")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("stats.json", jsonServers, 0644)
		if err != nil {
			panic(err)
		}
	}()

	// Start up server.
	r := router.Init()
	fmt.Println("Your Platytpus instance is live on port :9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}

// Generate barebones files required to run.
func genFiles() {
	if _, err := os.Stat(".secret"); os.IsNotExist(err) {
		fmt.Println("Generating secret key to .secret, use this to configure your servers")
		secretKey()
	}
	if _, err := os.Stat("stats.json"); os.IsNotExist(err) {
		err = ioutil.WriteFile("stats.json", []byte("{}"), 0644)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		err = ioutil.WriteFile("config.yml", []byte("port: 9090\ninterval: 5\n"), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func secretKey() {
	key, err := generateRandomString(32)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(".secret", []byte(key), 0644)
	if err != nil {
		panic(err)
	}
}

// From https://stackoverflow.com/questions/32349807/how-can-i-generate-a-random-int-using-the-crypto-rand-package
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return hex.EncodeToString(b), err
}
