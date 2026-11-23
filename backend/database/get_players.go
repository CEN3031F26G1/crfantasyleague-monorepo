
package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetPlayers retrieves all professional players from PostgreSQL.
func GetPlayers(ctx context.Context, db *pgxpool.Pool) ([]Player, error) {
	query := `
		SELECT
			id,
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
		FROM players
		ORDER BY name
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []Player

	for rows.Next() {
		var player Player

		err := rows.Scan(
			&player.ID,
			&player.ExternalID,
			&player.Name,
			&player.Country,
			&player.CountryCode,
			&player.Qualification,
			&player.Route,
			&player.Bio,
			&player.ImageURL,
			&player.SourceURL,
			&player.ImageSourceURL,
		)

		if err != nil {
			return nil, err
		}

		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}
