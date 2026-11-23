
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"crfantasyleague/backend/database"

	"github.com/joho/godotenv"
)

 
type CatalogPlayer struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Country       string `json:"country"`
	CountryCode   string `json:"countryCode"`
	Qualification string `json:"qualification"`
	Route         string `json:"route"`
	Bio           string `json:"bio"`
	Image         string `json:"image"`
	Source        string `json:"source"`
	ImageSource   string `json:"imageSource"`
}

func main() {
	//Load environment variables 
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load .env: ", err)
	}

	//Make sure the database URL exists
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Get the path to the players
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run ./cmd/importplayers <players.json>")
	}


	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Failed to read players.json: ", err)
	}

	//Convert JSON data 
	var players []CatalogPlayer
	if err := json.Unmarshal(data, &players); err != nil {
		log.Fatal("Failed to parse players.json: ", err)
	}

	if len(players) == 0 {
		log.Fatal("No players found in players.json")
	}


	ctx, cancel := context.WithTimeout(
		context.Background(),
		60*time.Second,
	)
	defer cancel()

	//connect to database
	db, err := database.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL: ", err)
	}
	defer db.Close()

	//import each player
	for _, p := range players {
		player := database.Player{
			ExternalID:     p.ID,
			Name:           p.Name,
			Country:        p.Country,
			CountryCode:    p.CountryCode,
			Qualification:  p.Qualification,
			Route:          p.Route,
			Bio:            p.Bio,
			ImageURL:       p.Image,
			SourceURL:      p.Source,
			ImageSourceURL: p.ImageSource,
		}

		// Save the player in database
		err := database.SavePlayer(ctx, db, player)
		if err != nil {
			log.Fatalf("Failed to save %s: %v", p.Name, err)
		}

		fmt.Println("Saved:", p.Name)
	}

	fmt.Printf("Successfully imported %d players!\n", len(players))
}
