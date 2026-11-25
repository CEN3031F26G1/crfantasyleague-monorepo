package database

import (
	"context"
	"testing"
)

func TestSaveMatchRejectsInvalidInput(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name  string
		match Match
	}{
		{"blank external ID", Match{ExternalID: "  ", TournamentName: "Test Tournament", Status: "completed"}},
		{"blank tournament name", Match{ExternalID: "test-match", TournamentName: "  ", Status: "completed"}},
		{"invalid status", Match{ExternalID: "test-match", TournamentName: "Test Tournament", Status: "unknown"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Invalid input should be rejected without using a database connection.
			if _, err := SaveMatch(ctx, nil, tc.match); err == nil {
				t.Fatal("SaveMatch accepted invalid input")
			}
		})
	}
}

func TestSavePlayerMatchStatsRejectsInvalidInput(t *testing.T) {
	ctx := context.Background()
	valid := PlayerMatchStats{MatchID: 1, PlayerID: 1, Result: "win"}
	cases := []struct {
		name   string
		mutate func(*PlayerMatchStats)
	}{
		{"missing match ID", func(s *PlayerMatchStats) { s.MatchID = 0 }},
		{"negative player ID", func(s *PlayerMatchStats) { s.PlayerID = -1 }},
		{"invalid result", func(s *PlayerMatchStats) { s.Result = "unknown" }},
		{"negative towers destroyed", func(s *PlayerMatchStats) { s.TowersDestroyed = -1 }},
		{"negative towers lost", func(s *PlayerMatchStats) { s.TowersLost = -1 }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stats := valid
			tc.mutate(&stats)
			if _, err := SavePlayerMatchStats(ctx, nil, stats); err == nil {
				t.Fatal("SavePlayerMatchStats accepted invalid input")
			}
		})
	}
}

func TestPlayerStatsQueriesRejectInvalidPlayerID(t *testing.T) {
	ctx := context.Background()
	for _, playerID := range []int64{0, -1} {
		if _, err := GetPlayerStats(ctx, nil, playerID); err == nil {
			t.Errorf("GetPlayerStats accepted player ID %d", playerID)
		}
		if _, err := GetPlayerTotals(ctx, nil, playerID); err == nil {
			t.Errorf("GetPlayerTotals accepted player ID %d", playerID)
		}
	}
}
