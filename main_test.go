package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

const APIBaseURL = "http://chuck.text"
const JokeText = "Don't worry about tests, Chuck Norris's test cases cover your code too."
const JokeCategory = "testing"

func TestSuccessfulRandomJoke(t *testing.T) {
	defer gock.Off()

	gock.New(APIBaseURL).
		Get("/jokes/random").
		MatchHeader("Accept", "text/plain").
		Reply(200).
		BodyString(JokeText)

	config := Config{
		APIBaseURL: APIBaseURL,
	}

	joke, err := getRandomJoke(config)

	require.NoError(t, err)
	require.Equal(t, JokeText, joke)
	require.True(t, gock.IsDone())
}

func TestSuccessfulCategorizedJoke(t *testing.T) {
	defer gock.Off()

	gock.New(APIBaseURL).
		Get("/jokes/random").
		MatchParam("category", JokeCategory).
		MatchHeader("Accept", "text/plain").
		Reply(200).
		BodyString(JokeText)

	config := Config{
		APIBaseURL: APIBaseURL,
		Category:   JokeCategory,
	}

	joke, err := getRandomJoke(config)

	require.NoError(t, err)
	require.Equal(t, JokeText, joke)
	require.True(t, gock.IsDone())
}

func TestUnsuccessfulHttpResponse(t *testing.T) {
	defer gock.Off()

	gock.New(APIBaseURL).
		Get("/jokes/random").
		Reply(500).
		BodyString("Internal server error")

	config := Config{
		APIBaseURL: APIBaseURL,
	}

	joke, err := getRandomJoke(config)

	require.Error(t, err, "Internal server error")
	require.Empty(t, joke)
	require.True(t, gock.IsDone())
}

func TestRequestError(t *testing.T) {
	defer gock.Off()

	gock.New(APIBaseURL).
		ReplyError(errors.New("test failure"))

	config := Config{
		APIBaseURL: APIBaseURL,
	}

	joke, err := getRandomJoke(config)

	require.Error(t, err, "test failure")
	require.Empty(t, joke)
	require.True(t, gock.IsDone())
}

func TestNoBaseApiUrl(t *testing.T) {
	joke, err := getRandomJoke(Config{})

	require.Error(t, err)
	require.Empty(t, joke)
}
