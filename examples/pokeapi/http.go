package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const fetchByNameURLMask = "https://pokeapi.co/api/v2/pokemon/%s" // (name)
var client *http.Client

func init() {
	client = http.DefaultClient
}

type pokemonTypeWithBindings struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type pokemonWithBindings struct {
	Types []pokemonTypeWithBindings `json:"types"`
}

// getPokemon fetches the pokemon from pokeapi using its name
func getPokemon(name string) (Pokemon, error) {
	resp, err := client.Get(fmt.Sprintf(fetchByNameURLMask, name))
	if err != nil {
		return Pokemon{}, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return Pokemon{}, errors.New("could not find pokemon")
	} else if resp.StatusCode >= http.StatusBadRequest {
		fmt.Println("error, got ", resp.StatusCode)
		return Pokemon{}, errors.New("failed to get pokemon")
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, nil
	}

	var p pokemonWithBindings
	if err := json.Unmarshal(bytes, &p); err != nil {
		return Pokemon{}, err
	}

	pokemon := Pokemon{
		Name:  name,
		Types: make([]PokemonType, len(p.Types)),
	}

	for _, t := range p.Types {
		pokemon.Types = append(pokemon.Types, PokemonType{
			Type: struct{ Name string }{t.Type.Name},
		})
	}

	return pokemon, nil
}
