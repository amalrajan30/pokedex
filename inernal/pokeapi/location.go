package pokeapi

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/amalrajan30/pokedexcli/inernal/pokecache"
)

func GetLocationArea(url string, cache *pokecache.Cache) (LocationArea, bool) {

	if cached, ok := cache.Get(url); ok {
		data, ok := PaseJSON[LocationArea](cached)
		if ok {
			return data, true
		}
	}

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	cache.Add(url, body)

	data, ok := PaseJSON[LocationArea](body)
	if ok {
		return data, true
	}

	return data, false
}


func ExplorePokemon(area string, cache *pokecache.Cache) (ExploreArea, bool) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area)

	if cached, ok := cache.Get(url); ok {
		data, ok := PaseJSON[ExploreArea](cached)
		if ok {
			return data, true
		}
	}

	body, err := MakeGetCall(url)

	data := ExploreArea{}

	if err != nil {
		return data, false
	}

	cache.Add(url, body)

	data, ok := PaseJSON[ExploreArea](body)

	if ok {
		return data, true
	}

	return data, false
}