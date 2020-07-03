package telegram

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
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
		ApiKey:   extractConfigValue(configStrings[0]),
		Channels: extractConfigValue(configStrings[1]),
	}

	requestData := fmt.Sprintf(`?chat_id=%v&text=%v`,
		config.Channels, url.QueryEscape(data+" just went offline!"))
	requestUrl := "https://api.telegram.org/bot" + config.ApiKey

	request, err := http.NewRequest("GET", requestUrl+"/sendMessage"+requestData, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return "", nil
}

// Dirty way of extracting the values.
func extractConfigValue(line string) string {
	return strings.SplitN(strings.Replace(line, " ", "", -1), ":", 2)[1]
}
