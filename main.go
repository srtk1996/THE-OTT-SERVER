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

// Movie struct definition
type Movie struct {
	ID       string     `json:"id"`
	Isbn     string     `json:"isbn"`
	Title    string     `json:"title"`
	Director *Director  `json:"director"`
}

// Director struct definition
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// In-memory database (slice of movies)
var movies []Movie

func main() {
	// Initialize Mux Router
	r := mux.NewRouter()

	// Add some initial movies to our in-memory database
	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "438227",
		Title: "Movie One",
		Director: &Director{
			Firstname: "John",
			Lastname:  "Doe",
		},
	})
	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "454555",
		Title: "Movie Two",
		Director: &Director{
			Firstname: "Steve",
			Lastname:  "Smith",
		},
	})

	// Define HTTP routes
	r.HandleFunc("/movies", getMovies).Methods("GET")         // Get all movies
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")     // Get a movie by its ID
	r.HandleFunc("/movies", createMovie).Methods("POST")      // Create a new movie
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")  // Update a movie by its ID
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE") // Delete a movie by its ID

	// Start the HTTP server
	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}

// Handler to get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set response header to JSON
	json.NewEncoder(w).Encode(movies)
// 	The json.NewEncoder() function creates a new JSON encoder that writes its output directly to the provided writer.
// In this case, the writer is w, which is the HTTP response writer (http.ResponseWriter).
// The .Encode() method takes the movies variable and converts it into JSON.
// Once encoded, the JSON data is automatically written to the ResponseWriter.
}

// Handler to get a single movie by ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get parameters from the URL
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}

// Handler to create a new movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000)) // Generate a random ID for the new movie
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

// Handler to update an existing movie by ID
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...) // Remove the old movie
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"] // Ensure the movie retains its original ID
			movies = append(movies, movie) // Add the updated movie
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}

// Handler to delete a movie by ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			fmt.Fprintf(w, "Movie with ID %s has been deleted", params["id"])
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}