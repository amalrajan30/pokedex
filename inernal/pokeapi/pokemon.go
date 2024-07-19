package pokeapi

import (
	"fmt"

	"github.com/amalrajan30/pokedexcli/inernal/pokecache"
)

func GetPokemon(name string, cache *pokecache.Cache) (Pokemon, bool) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	if cached, ok := cache.Get(url); ok {
		data, ok := PaseJSON[Pokemon](cached)
		if ok {
			return data, true
		}
	}

	body, err := MakeGetCall(url)

	data := Pokemon{}

	if err != nil {
		return data, false
	}

	cache.Add(url, body)

	data, ok := PaseJSON[Pokemon](body)

	if ok {
		return data, true
	}

	return data, false
}