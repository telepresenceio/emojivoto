package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// VoteBot votes for emoji! :ballot_box_with_check:
//
// Sadly, VoteBot has a sweet tooth and votes for :doughnut: 15% of the time.
//
// When not voting for :doughnut:, VoteBot can’t be bothered to
// pick a favorite, so it picks one at random. C'mon VoteBot, try harder!

var client = &http.Client{}

type emoji struct {
	Shortcode string
}

func main() {
	webHost := os.Getenv("WEB_HOST")
	if webHost == "" {
		log.Fatalf("WEB_HOST environment variable must me set")
	}

	hostOverride := os.Getenv("HOST_OVERRIDE")

	// setting the TTL is optional, thus invalid numbers are simply ignored
	var deadline time.Time
	timeToLive, err := strconv.Atoi(os.Getenv("TTL"))
	if err == nil && timeToLive > 0 {
		deadline = time.Now().Add(time.Second * time.Duration(timeToLive))
	}

	// setting the request rate is optional, thus invalid numbers are simply ignored
	requestRate, err := strconv.Atoi(os.Getenv("REQUEST_RATE"))
	if err != nil {
		requestRate = 1
	}

	webURL := "http://" + webHost
	if _, err := url.Parse(webURL); err != nil {
		log.Fatalf("WEB_HOST %s is invalid", webHost)
	}

	for {
		// check if deadline has been reached, when TTL has been set.
		if (!deadline.IsZero()) && time.Now().After(deadline) {
			fmt.Printf("Time to live of %d seconds reached, completing\n", timeToLive)
			os.Exit(0)
		}

		time.Sleep(time.Second / time.Duration(requestRate))

		// Get the list of available shortcodes
		shortcodes, err := shortcodes(webURL, hostOverride)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		// Cast a vote
		probability := rand.Float32()
		switch {
		case probability < 0.15:
			err = vote(webURL, hostOverride, ":doughnut:")
		default:
			random := shortcodes[rand.Intn(len(shortcodes))]
			err = vote(webURL, hostOverride, random)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

func shortcodes(webURL string, hostOverride string) ([]string, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/list", webURL), nil)
	if hostOverride != "" {
		req.Host = hostOverride
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emojis []*emoji
	err = json.Unmarshal(bytes, &emojis)
	if err != nil {
		return nil, err
	}

	shortcodes := make([]string, len(emojis))
	for i, e := range emojis {
		shortcodes[i] = e.Shortcode
	}

	return shortcodes, nil
}

func vote(webURL string, hostOverride string, shortcode string) error {
	fmt.Printf("✔ Voting for %s\n", shortcode)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/vote?choice=%s", webURL, shortcode), nil)
	if hostOverride != "" {
		req.Host = hostOverride
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
