package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestSavePlayer(t *testing.T) {
	// Load PostgreSQL configuration.
	_ = godotenv.Load("../.env")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not configured")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	// Connect to PostgreSQL.
	db, err := Connect(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a unique test player ID.
	testID := fmt.Sprintf("test-%d", time.Now().UnixNano())

	// Delete the test player when the test finishes.
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(), 5*time.Second,
		)
		defer cleanupCancel()

		_, err := db.Exec(
			cleanupCtx,
			"DELETE FROM players WHERE external_id = $1",
			testID,
		)
		if err != nil {
			t.Errorf("Cleanup failed: %v", err)
		}
	}()

	// =====================================
	// TEST 1: INSERT PLAYER
	// =====================================

	player := Player{
		ExternalID:  testID,
		Name:        "Test Player",
		Country:     "Japan",
		CountryCode: "JP",
	}

	// Save the new player.
	if err := SavePlayer(ctx, db, player); err != nil {
		t.Fatal("Failed to save player:", err)
	}

	// Retrieve the player from PostgreSQL.
	var name string

	err = db.QueryRow(
		ctx,
		"SELECT name FROM players WHERE external_id = $1",
		testID,
	).Scan(&name)

	if err != nil {
		t.Fatal("Failed to retrieve player:", err)
	}

	if name != "Test Player" {
		t.Fatalf("Unexpected player name: %s", name)
	}

	t.Log("Player successfully saved and retrieved!")

	// =====================================
	// TEST 2: UPDATE PLAYER
	// =====================================

	// Change the existing player's name.
	player.Name = "Updated Test Player"

	// Save again using the same external ID.
	if err := SavePlayer(ctx, db, player); err != nil {
		t.Fatal("Failed to update player:", err)
	}

	// Verify the update and check for duplicates.
	var updatedName string
	var count int

	err = db.QueryRow(ctx, `
        SELECT MAX(name), COUNT(*)
        FROM players
        WHERE external_id = $1
    `, testID).Scan(&updatedName, &count)

	if err != nil {
		t.Fatal("Failed to retrieve updated player:", err)
	}

	if updatedName != "Updated Test Player" {
		t.Fatalf(
			"Expected updated name, got: %s",
			updatedName,
		)
	}

	if count != 1 {
		t.Fatalf(
			"Expected 1 player record, got: %d",
			count,
		)
	}

	t.Log("Player successfully updated without duplicates!")
}
