package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type PokemonType struct {
	Type struct {
		Name string
	}
}

type Pokemon struct {
	Name  string
	Types []PokemonType
}

type FetchResult struct {
	Pokemon
	Err error
}

func readPokemonNames() []string {
	f, err := os.Open("pokemons.txt")

	if err != nil {
		log.Fatalf("could not open pokemons.txt: %s", err.Error())
	}
	defer f.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(f)

	names := make([]string, 0, 1000)
	for scanner.Scan() {
		text := scanner.Text()
		names = append(names, strings.TrimSpace(text))
	}

	return names
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("elapsed time: %vms\n", time.Since(start).Milliseconds())
	}()

	names := readPokemonNames()

	jobs := make(chan string)   // contains URLs to be fetched
	done := make(chan struct{}) // when closed, indicates work has finished

	results := make(chan FetchResult) // results to be collected, embed of pokemon (name and type) + error if sth happened

	// producer
	go func() {
		for _, pokemonName := range names {
			jobs <- pokemonName
		}
	}()

	// 50 workers
	for workers := 0; workers < 50; workers++ {
		go func() {
			for {
				select {
				case <-time.After(1 * time.Minute): // timeout
					close(done)
					return
				case pokemon := <-jobs:
					obj, err := getPokemon(pokemon)
					results <- FetchResult{
						Pokemon: Pokemon{
							Name:  pokemon,
							Types: obj.Types,
						},
						Err: err,
					}
				case <-done:
					return
				}
			}
		}()
	}

	countByType := make(map[string]int)
	errorNames := make([]string, 0)

	// consumer (collects results)
	for range names {
		fetchResult := <-results
		if fetchResult.Err != nil {
			errorNames = append(errorNames, fetchResult.Name)
		} else {
			for _, pokemonType := range fetchResult.Types {
				countByType[pokemonType.Type.Name]++
			}
		}
	}

	close(done)

	fmt.Println("errors: ", errorNames)
	fmt.Println("count by type: ", countByType)
}
