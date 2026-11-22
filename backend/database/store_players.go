package database

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Player represents a professional Clash Royale player.
type Player struct {
	ID             int64
	ExternalID     string
	Name           string
	Country        string
	CountryCode    string
	Qualification  string
	Route          string
	Bio            string
	ImageURL       string
	SourceURL      string
	ImageSourceURL string
}

// SavePlayer inserts or updates a professional player.
func SavePlayer(
	ctx context.Context,
	db *pgxpool.Pool,
	player Player,
) error {

	if strings.TrimSpace(player.ExternalID) == "" ||
		strings.TrimSpace(player.Name) == "" {
		return errors.New("player ID and name are required")
	}

	query := `
        INSERT INTO players (
            external_id,
            name,
            country,
            country_code,
            qualification,
            route,
            bio,
            image_url,
            source_url,
            image_source_url
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (external_id)
        DO UPDATE SET
            name = EXCLUDED.name,
            country = EXCLUDED.country,
            country_code = EXCLUDED.country_code,
            qualification = EXCLUDED.qualification,
            route = EXCLUDED.route,
            bio = EXCLUDED.bio,
            image_url = EXCLUDED.image_url,
            source_url = EXCLUDED.source_url,
            image_source_url = EXCLUDED.image_source_url,
            updated_at = CURRENT_TIMESTAMP
    `

	_, err := db.Exec(
		ctx,
		query,
		player.ExternalID,
		player.Name,
		player.Country,
		player.CountryCode,
		player.Qualification,
		player.Route,
		player.Bio,
		player.ImageURL,
		player.SourceURL,
		player.ImageSourceURL,
	)

	return err
}
