package database

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)
//one match
type Match struct {
	ID 	  int64
	ExternalID string
	TournamentName string
	MatchDate *time.Time
	Status string
}
//inserts match
func SaveMatch(ctx context.Context, db *pgxpool.Pool, match Match) (int64, error) {
	if strings.TrimSpace(match.ExternalID) == "" || strings.TrimSpace(match.TournamentName) == "" {
		return 0, errors.New("match ID and tournament name cannot be empty")
	}

	if match.Status != "scheduled" && match.Status != "ongoing" && match.Status != "completed" {
		return 0, errors.New("invalid match status")
	}

	//updates exisiting match
	query := `
		INSERT INTO matches (external_id, tournament_name, match_date, status)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (external_id) DO UPDATE 
		SET tournament_name = EXCLUDED.tournament_name,
		    match_date = EXCLUDED.match_date,
		    status = EXCLUDED.status,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id;
	`

	var matchID int64
	err := db.QueryRow(ctx, query, match.ExternalID, match.TournamentName, match.MatchDate, match.Status).Scan(&matchID)

	if err != nil {
		return 0, err
	}
	return matchID, nil
}
