package database 

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)
//players performance in a match
type PlayerMatchStats struct {
	ID 	  int64
	MatchID int64
	PlayerID int64
	Result string
	TowersDestroyed int
	TowersLost int
}

//inserts or updates player match stats
func SavePlayerMatchStats(ctx context.Context, db *pgxpool.Pool, stats PlayerMatchStats) (int64, error) {
	if stats.MatchID <= 0 || stats.PlayerID <= 0 {
		return 0, errors.New("match ID and player ID are required")
	}

	if stats.Result != "win" && stats.Result != "loss" && stats.Result != "draw" {
		return 0, errors.New("invalid match result")
	}

	if stats.TowersDestroyed < 0 || stats.TowersLost < 0 {
		return 0, errors.New("tower counts cannot be negative")
	}

	query := `
		INSERT INTO player_match_stats (match_id, player_id, result, towers_destroyed, towers_lost)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (match_id, player_id) DO UPDATE 
		SET result = EXCLUDED.result,
		    towers_destroyed = EXCLUDED.towers_destroyed,
		    towers_lost = EXCLUDED.towers_lost,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id;
	`
	
	var statsID int64
	err := db.QueryRow(ctx, query, stats.MatchID, stats.PlayerID, stats.Result, stats.TowersDestroyed, stats.TowersLost).Scan(&statsID)

	if err != nil {
		return 0, err
	}
	return statsID, nil
}
