package repository

import (
	"errors"
	"testing"

	"github.com/orienteering/platform/internal/model"
)

func TestUserRepository_CreateAndGetByUsername(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{Username: "alice", PasswordHash: "hash", Nickname: "爱丽丝", Role: "user", Points: 100}
	if err := repo.Create(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	got, err := repo.GetByUsername("alice")
	if err != nil {
		t.Fatalf("get user by username: %v", err)
	}
	if got.ID != user.ID || got.Nickname != "爱丽丝" {
		t.Fatalf("unexpected user: %+v", got)
	}
	if _, err := repo.GetByUsername("nobody"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	_ = repo.Create(&model.User{Username: "bob", PasswordHash: "h", Role: "user"})

	cases := []struct {
		name     string
		username string
		want     bool
	}{
		{"existing", "bob", true},
		{"missing", "charlie", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.ExistsByUsername(tc.username)
			if err != nil {
				t.Fatalf("exists by username: %v", err)
			}
			if got != tc.want {
				t.Fatalf("ExistsByUsername(%q) = %v, want %v", tc.username, got, tc.want)
			}
		})
	}
}

func TestUserRepository_AddPoints(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	user := &model.User{Username: "carol", PasswordHash: "h", Role: "user", Points: 50}
	_ = repo.Create(user)

	if err := repo.AddPoints(user.ID, 30); err != nil {
		t.Fatalf("add points: %v", err)
	}
	got, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if got.Points != 80 {
		t.Fatalf("points = %d, want 80", got.Points)
	}
}

func TestUserRepository_List(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	for _, name := range []string{"alice", "bob", "carol"} {
		_ = repo.Create(&model.User{Username: name, Nickname: name, Role: "user"})
	}
	users, total, err := repo.List(1, 2, 0, "alice")
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if total != 1 || len(users) != 1 || users[0].Username != "alice" {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(users))
	}
}
