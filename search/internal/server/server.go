package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tprei/semcomp/search/internal/handler"
)

func setup(w http.ResponseWriter) {
	// inform writer that its writing json
	w.Header().Set("Content-Type", "application/json")

	// allow all origins
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func Search() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { // default server panics to status 500
			if r := recover(); r != nil {
				err := r.(error)
				fmt.Println("an error occured:", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()

		setup(w)

		queryValues := r.URL.Query()
		query := queryValues.Get("query")

		if r.URL.Path != "/search" || len(query) == 0 { // wrong path or empty query
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(`[]`)
			return
		}

		songs, err := handler.Search(query)
		if err != nil {
			panic(fmt.Errorf("error searching for query %s: %s", query, err))
		} else { // marshal results to json
			if bytes, err := json.MarshalIndent(songs, "", "    "); err != nil {
				panic("error marshalling results")
			} else {
				w.Write(bytes)
				return
			}
		}
	}
}

func RunServer() error {
	addr := "localhost"
	port := "8080"

	fmt.Printf("running server on %s:%s\n", addr, port)
	return http.ListenAndServe(
		fmt.Sprintf("%s:%s", addr, port),
		http.HandlerFunc(Search()),
	)
}
