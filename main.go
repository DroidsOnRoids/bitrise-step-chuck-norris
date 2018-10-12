package main

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Config ...
type Config struct {
	APIBaseURL string `env:"api_base_url,required"`
	Category   string `env:"category"`
}

func main() {
	var config Config
	if err := stepconf.Parse(&config); err != nil {
		log.Errorf("Could not validate config, error: %s\n", err)
		os.Exit(1)
	}

	joke, err := getRandomJoke(config)

	if err != nil {
		log.Errorf("Could not get joke, error: %s\n", err)
		os.Exit(1)
	}

	log.Donef("%s", joke)

	if err := exportEnvironmentWithEnvman(joke); err != nil {
		log.Errorf("Could not export joke, error: %s\n", err)
		os.Exit(1)
	}
}

func getRandomJoke(config Config) (string, error) {
	request, err := buildJokeRequest(config)
	if err != nil {
		return "", err
	}

	response, err := getJoke(request)
	if err != nil {
		return "", err
	}

	defer func() {
		if response.Body.Close() != nil {
			log.Warnf("Failed to close response body, error: %s", response.Body.Close())
		}
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned an error: %s", response.Status)
	}

	content, err := ioutil.ReadAll(response.Body)
	return string(content), err
}

func buildJokeRequest(config Config) (*http.Request, error) {
	jokeURL, err := buildJokeURL(config)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", jokeURL.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "text/plain")
	return request, nil
}

func getJoke(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	client.Timeout = 20 * time.Second
	return client.Do(request)
}

func buildJokeURL(config Config) (*url.URL, error) {
	jokeURL, err := url.Parse(config.APIBaseURL)

	if err != nil {
		return nil, err
	}

	jokeURL.Path = "/jokes/random"

	if len(config.Category) > 0 {
		query := jokeURL.Query()
		query.Set("category", config.Category)
		jokeURL.RawQuery = query.Encode()
	}

	return jokeURL, nil
}

func exportEnvironmentWithEnvman(value string) error {
	return command.RunCommand("envman", "add", "--key", "CHUCK_NORRIS_JOKE", "--value", value)
}