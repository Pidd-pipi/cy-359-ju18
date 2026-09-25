package service

import (
	"testing"

	"github.com/orienteering/platform/internal/dto"
)

func TestSortFinished(t *testing.T) {
	rows := []dto.LeaderboardRow{
		{TeamID: 3, TotalSeconds: 900},
		{TeamID: 1, TotalSeconds: 600},
		{TeamID: 2, TotalSeconds: 900},
	}
	got := sortFinished(rows)
	want := []int64{1, 2, 3}
	for i, id := range want {
		if got[i].TeamID != id {
			t.Fatalf("position %d team = %d, want %d", i, got[i].TeamID, id)
		}
	}
}

func TestSortInProgress(t *testing.T) {
	rows := []dto.LeaderboardRow{
		{TeamID: 5, CheckpointCount: 1},
		{TeamID: 6, CheckpointCount: 3},
		{TeamID: 7, CheckpointCount: 3},
	}
	got := sortInProgress(rows)
	want := []int64{6, 7, 5}
	for i, id := range want {
		if got[i].TeamID != id {
			t.Fatalf("position %d team = %d, want %d", i, got[i].TeamID, id)
		}
	}
}
