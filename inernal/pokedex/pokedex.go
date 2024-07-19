package pokedex

import "github.com/amalrajan30/pokedexcli/inernal/pokeapi"

type PokeDex struct {
	Dex map[string]pokeapi.Pokemon
}

func(d PokeDex) AddPokemon(pokemon pokeapi.Pokemon, name string) {
	d.Dex[name] = pokemon
}

func NewDex() PokeDex {
	newDex := PokeDex{
		Dex: make(map[string]pokeapi.Pokemon),
	}

	return newDex
}