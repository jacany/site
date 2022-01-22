package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shkh/lastfm-go/lastfm"
)

var cachedAlbums string

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hello!")
    fmt.Println("Endpoint Hit: homePage")
}

func topAlbums(w http.ResponseWriter, r *http.Request){
	setupCorsResponse(&w, r)

    fmt.Fprintf(w, cachedAlbums)
    fmt.Println("Endpoint Hit: topAlbums")
}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func main() {
	fmt.Println("starting...")
	// load env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s\n", err)
	}
	ApiKey := os.Getenv("API_KEY")
	Port := os.Getenv("PORT")

	// lastfm stuff
	api := lastfm.New(ApiKey, "")
	go func() {
		for {
			result, err := api.User.GetTopAlbums(lastfm.P{"user": "realshapes"}) //discarding error
			if err != nil {
				log.Printf("error getting top albums: %s\n", err)
			}
			
			marshalled, err := json.Marshal(result)
			if err != nil {
				log.Printf("error marshalling top albums: %s\n", err)
			}
			cachedAlbums = string(marshalled)
			fmt.Println(cachedAlbums)
			fmt.Println("(re)cached top albums")
			time.Sleep(86400 * time.Second)
		}
	}()

	// start web server
    http.HandleFunc("/", homePage)
	http.HandleFunc("/topalbums", topAlbums)
	log.Printf("listening on port %s\n", Port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port), nil))
}