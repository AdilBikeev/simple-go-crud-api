package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json: "id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {

		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.Id = strconv.Itoa(rand.Intn(1000000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movieUpdate Movie

	for i, movie := range movies {
		if params["id"] == movie.Id {
			// Delete movie
			movies = append(movies[:i], movies[i+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&movieUpdate)

			// Add movie
			movieUpdate.Id = params["id"]
			movies = append(movies, movieUpdate)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{message}: "Not found movie"`))
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{
		Id:    "1",
		Isbn:  "438277",
		Title: "Movie One",
		Director: &Director{
			Firstname: "Ivan",
			Lastname:  "Ivanov",
		}})
	movies = append(movies, Movie{
		Id:    "2",
		Isbn:  "45446",
		Title: "Movie Two",
		Director: &Director{
			Firstname: "Petr",
			Lastname:  "Petrov",
		}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
