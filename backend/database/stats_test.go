package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
)

func TestSaveMatchAndPlayerStats(t *testing.T) {
	_ = godotenv.Load("../.env")
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := Connect(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testID := fmt.Sprintf("pbi06-test-%d", time.Now().UnixNano())
	secondMatchExternalID := testID + "-second"

	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()

		if _, err := db.Exec(
			cleanupCtx,
			"DELETE FROM matches WHERE external_id IN ($1, $2)",
			testID,
			secondMatchExternalID,
		); err != nil {
			t.Errorf("failed to clean up matches: %v", err)
		}

		if _, err := db.Exec(
			cleanupCtx,
			"DELETE FROM players WHERE external_id = $1",
			testID,
		); err != nil {
			t.Errorf("failed to clean up player: %v", err)
		}
	}()

	if err := SavePlayer(ctx, db, Player{
		ExternalID: testID,
		Name:       "Temporary Test Player",
	}); err != nil {
		t.Fatal(err)
	}

	var playerID int64
	if err := db.QueryRow(ctx,
		"SELECT id FROM players WHERE external_id = $1",
		testID,
	).Scan(&playerID); err != nil {
		t.Fatal(err)
	}

	
	matchID, err := SaveMatch(ctx, db, Match{
		ExternalID:     testID,
		TournamentName: "Test Tournament",
		Status:         "completed",
	})
	if err != nil {
		t.Fatal(err)
	}

	const updatedTournamentName = "Updated Test Tournament"
	updatedMatchID, err := SaveMatch(ctx, db, Match{
		ExternalID:     testID,
		TournamentName: updatedTournamentName,
		Status:         "completed",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedMatchID != matchID {
		t.Fatalf("match update returned ID %d, want %d", updatedMatchID, matchID)
	}

	var matchCount int
	var tournamentName string
	if err := db.QueryRow(ctx, `
         SELECT COUNT(*), COALESCE(MAX(tournament_name), '')
         FROM matches
         WHERE external_id = $1
      `, testID).Scan(&matchCount, &tournamentName); err != nil {
		t.Fatal(err)
	}
	if matchCount != 1 || tournamentName != updatedTournamentName {
		t.Fatalf("expected one match named %q, got %d rows with name %q",
			updatedTournamentName, matchCount, tournamentName)
	}

	// Use a positive ID that does not exist in the test database.
	// Zero or negative IDs would be rejected by Go before reaching PostgreSQL.
	const missingID int64 = 1<<63 - 1

	cases := []struct {
		name       string
		matchID    int64
		playerID   int64
		constraint string
	}{
		{
			name:       "nonexistent player",
			matchID:    matchID,
			playerID:   missingID,
			constraint: "player_match_stats_player_id_fkey",
		},
		{
			name:       "nonexistent match",
			matchID:    missingID,
			playerID:   playerID,
			constraint: "player_match_stats_match_id_fkey",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := SavePlayerMatchStats(ctx, db, PlayerMatchStats{
				MatchID:  tc.matchID,
				PlayerID: tc.playerID,
				Result:   "win",
			})

			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				t.Fatalf("expected a PostgreSQL foreign key error, got %v", err)
			}
			if pgErr.Code != "23503" || pgErr.ConstraintName != tc.constraint {
				t.Fatalf("got code %s on constraint %q; want 23503 on %q",
					pgErr.Code, pgErr.ConstraintName, tc.constraint)
			}
		})
	}

	// Save stats, then update them for the same player and match.
	statsID, err := SavePlayerMatchStats(ctx, db, PlayerMatchStats{
		MatchID:         matchID,
		PlayerID:        playerID,
		Result:          "win",
		TowersDestroyed: 2,
		TowersLost:      1,
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedStatsID, err := SavePlayerMatchStats(ctx, db, PlayerMatchStats{
		MatchID:         matchID,
		PlayerID:        playerID,
		Result:          "win",
		TowersDestroyed: 3,
		TowersLost:      0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedStatsID != statsID {
		t.Fatalf("stats update returned ID %d, want %d", updatedStatsID, statsID)
	}

	// Verify what PostgreSQL actually stored.
	var result string
	var towersDestroyed, towersLost int
	if err := db.QueryRow(ctx, `
         SELECT result, towers_destroyed, towers_lost
         FROM player_match_stats
         WHERE match_id = $1 AND player_id = $2
      `, matchID, playerID).Scan(&result, &towersDestroyed, &towersLost); err != nil {
		t.Fatal(err)
	}
	if result != "win" || towersDestroyed != 3 || towersLost != 0 {
		t.Fatalf("unexpected stored stats: result=%q, destroyed=%d, lost=%d",
			result, towersDestroyed, towersLost)
	}

	savedStats, err := GetPlayerStats(ctx, db, playerID)
	if err != nil {
		t.Fatal(err)
	}
	if len(savedStats) != 1 {
		t.Fatalf("expected 1 stats record, got %d", len(savedStats))
	}

	wantStats := PlayerMatchStats{
		ID:              statsID,
		MatchID:         matchID,
		PlayerID:        playerID,
		Result:          "win",
		TowersDestroyed: 3,
		TowersLost:      0,
	}
	if savedStats[0] != wantStats {
		t.Fatalf("got stats %+v, want %+v", savedStats[0], wantStats)
	}

	totals, err := GetPlayerTotals(ctx, db, playerID)
	if err != nil {
		t.Fatal(err)
	}
	wantTotals := PlayerTotals{
		PlayerID:        playerID,
		TotalMatches:    1,
		Wins:            1,
		TowersDestroyed: 3,
	}
	if totals != wantTotals {
		t.Fatalf("got single-match totals %+v, want %+v", totals, wantTotals)
	}

	secondMatchID, err := SaveMatch(ctx, db, Match{
		ExternalID:     secondMatchExternalID,
		TournamentName: "Test Tournament",
		Status:         "completed",
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := SavePlayerMatchStats(ctx, db, PlayerMatchStats{
		MatchID:         secondMatchID,
		PlayerID:        playerID,
		Result:          "loss",
		TowersDestroyed: 1,
		TowersLost:      2,
	}); err != nil {
		t.Fatal(err)
	}

	totals, err = GetPlayerTotals(ctx, db, playerID)
	if err != nil {
		t.Fatal(err)
	}
	wantTotals = PlayerTotals{
		PlayerID:        playerID,
		TotalMatches:    2,
		Wins:            1,
		Losses:          1,
		TowersDestroyed: 4,
		TowersLost:      2,
	}
	if totals != wantTotals {
		t.Fatalf("got two-match totals %+v, want %+v", totals, wantTotals)
	}
}
