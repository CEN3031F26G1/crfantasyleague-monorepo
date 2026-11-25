package database 

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

//retrieve match stats
func GetPlayerStats(ctx context.Context, db *pgxpool.Pool, playerID int64) ([]PlayerMatchStats, error) {
	if playerID <= 0 {
		return nil, errors.New("invalid player ID")
	}

	//retrive stats only for player provided
	query := `
		SELECT id, match_id, player_id, result, towers_destroyed, towers_lost
		FROM player_match_stats
		WHERE player_id = $1
		ORDER BY match_id DESC
	`

	rows, err := db.Query(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []PlayerMatchStats{}
	for rows.Next(){
		var s PlayerMatchStats
		err := rows.Scan(&s.ID, &s.MatchID, &s.PlayerID, &s.Result, &s.TowersDestroyed, &s.TowersLost)

		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}