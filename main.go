package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Rabbidking/pokedex/internal/pokecache"
)

const baseLocationAreaURL = "https://pokeapi.co/api/v2/location-area/"

var cache = pokecache.NewCache(time.Duration(10 * time.Second))

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	//store next and previous URLs
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type pokeAPIResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

var commandRegistry = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Displays the names of the next 20 locations",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays the names of the previous 20 locations",
		callback:    commandMapb,
	},
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Print("\n")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: display the names of the next 20 locations")
	fmt.Println("mapb: display the previous 20 locations")
	return nil
}

func fetchAndPrintLocations(url string, cfg *config) error {
	//Fetch data from URL, unmarshal the JSON response into our struct, print each result, update cfg.Next and cfg.Previous, and return any errors
	var pokeResp pokeAPIResponse
	var data []byte
	var err error

	//check if data from url is already in our cache
	data, ok := cache.Get(url)
	if !ok {
		//fetch from http.Get, since it's not in our cache
		resp, err := http.Get(url)

		if err != nil {
			return err
		}

		//read body as []byte
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		data = body
		cache.Add(url, data)

	}

	//Unmarshal the JSON into our struct
	err = json.Unmarshal(data, &pokeResp)
	if err != nil {
		return err
	}

	//print names of locations
	for _, field := range pokeResp.Results {
		fmt.Println(field.Name)
	}

	//update cfg.Next and cfg.Previous with the new values from pokeResp
	cfg.Next = pokeResp.Next
	cfg.Previous = pokeResp.Previous

	return nil
}

func commandMap(cfg *config) error {
	var urlToFetch string

	//determine which URL to fetch, per map call
	if cfg.Next == nil {
		//if nil, we start from the first page
		urlToFetch = baseLocationAreaURL
	} else {
		//we are good to use the next URL
		urlToFetch = *cfg.Next
	}

	return fetchAndPrintLocations(urlToFetch, cfg)
}

func commandMapb(cfg *config) error {

	//go back to previous page. Check if we're already on the first page!
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	return fetchAndPrintLocations(*cfg.Previous, cfg)
}

func cleanInput(text string) []string {
	output := []string{}
	split_text := strings.Fields(text)
	for _, field := range split_text {
		output = append(output, strings.ToLower(field))
	}
	return output
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		s := scanner.Text()
		clean_string := cleanInput(s)

		if len(clean_string) == 0 {
			// continue to next loop iteration
			continue
		}

		cmd := clean_string[0]
		command, ok := commandRegistry[cmd]
		if ok {
			command.callback(cfg)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
