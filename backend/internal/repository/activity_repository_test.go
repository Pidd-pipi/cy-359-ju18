package repository

import (
	"testing"
	"time"

	"github.com/orienteering/platform/internal/model"
)

func TestActivityRepository_CreateGetAndList(t *testing.T) {
	db := newTestDB(t)
	repo := NewActivityRepository(db)

	now := time.Now()
	act := &model.Activity{
		Title: "滨江夜跑定向赛", Difficulty: "adult", Status: "draft",
		CreatorID: 1, StartTime: now, EndTime: now.Add(3 * time.Hour),
	}
	if err := repo.Create(act); err != nil {
		t.Fatalf("create activity: %v", err)
	}
	got, err := repo.GetByID(act.ID)
	if err != nil {
		t.Fatalf("get activity: %v", err)
	}
	if got.Title != "滨江夜跑定向赛" {
		t.Fatalf("unexpected activity: %+v", got)
	}

	cases := []struct {
		name       string
		status     string
		difficulty string
		keyword    string
		wantTotal  int64
	}{
		{"all", "", "", "", 1},
		{"by status draft", "draft", "", "", 1},
		{"by status published", "published", "", "", 0},
		{"by difficulty adult", "", "adult", "", 1},
		{"by keyword", "", "", "滨江", 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, total, err := repo.List(1, 10, 0, tc.status, tc.difficulty, tc.keyword)
			if err != nil {
				t.Fatalf("list activities: %v", err)
			}
			if total != tc.wantTotal {
				t.Fatalf("total = %d, want %d", total, tc.wantTotal)
			}
		})
	}
}

func TestActivityRepository_CountCheckpoints(t *testing.T) {
	db := newTestDB(t)
	actRepo := NewActivityRepository(db)
	cpRepo := NewCheckpointRepository(db)

	act := &model.Activity{Title: "测试线路", Status: "draft", CreatorID: 1}
	_ = actRepo.Create(act)
	_ = cpRepo.Create(&model.Checkpoint{ActivityID: act.ID, Name: "CP1", Sequence: 1, QRCode: "qr1"})
	_ = cpRepo.Create(&model.Checkpoint{ActivityID: act.ID, Name: "CP2", Sequence: 2, QRCode: "qr2"})

	count, err := actRepo.CountCheckpoints(act.ID)
	if err != nil {
		t.Fatalf("count checkpoints: %v", err)
	}
	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
}
