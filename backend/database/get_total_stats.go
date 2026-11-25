package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

//total overall match stats for a player
type PlayerTotals struct {
	PlayerID int64
	TotalMatches int64
	Wins int64
	Losses int64
	Draws int64
	TowersDestroyed int64
	TowersLost int64
}

//players overall stats
func GetPlayerTotals(ctx context.Context, db *pgxpool.Pool, playerID int64) (PlayerTotals, error) {
	if playerID <= 0 {
		return PlayerTotals{}, errors.New("invalid player ID")

	}

	query := `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE result = 'win'),
			COUNT(*) FILTER (WHERE result = 'loss'),
			COUNT(*) FILTER (WHERE result = 'draw'),
			COALESCE(SUM(towers_destroyed), 0),
			COALESCE(SUM(towers_lost), 0)
		FROM player_match_stats
		WHERE player_id = $1
	`

	totals := PlayerTotals{PlayerID: playerID}
	err := db.QueryRow(ctx, query, playerID).Scan(
		&totals.TotalMatches,
		&totals.Wins,
		&totals.Losses,
		&totals.Draws,
		&totals.TowersDestroyed,
		&totals.TowersLost,
	)

	if err != nil {
		return PlayerTotals{}, err
	}

	return totals, nil
}