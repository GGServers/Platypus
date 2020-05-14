package telegram

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type TelegramConfiguration struct {
	ApiKey   string `yaml:"apikey"`
	Channels string `yaml:"channels"`
}

func Offline(data string) (string, error) {
	file, err := os.Open("plugins/telegram/config.yml")
	if err != nil {
		return "", err
	}

	var configStrings []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		configStrings = append(configStrings, scanner.Text())
	}
	config := TelegramConfiguration{
		ApiKey:   configStrings[0],
		Channels: configStrings[1],
	}

	requestData := fmt.Sprintf(`?chat_id=%v&text=%v`,
		config.Channels, url.QueryEscape(data + " just went offline!"))
	requestUrl := "https://api.telegram.org/bot" + config.ApiKey

	request, err := http.NewRequest("GET", requestUrl + "/sendMessage" + requestData, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	return "", nil
}
