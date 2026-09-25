package util

import (
	"testing"

	"github.com/orienteering/platform/internal/constants"
)

func TestActivityStatusText(t *testing.T) {
	cases := []struct {
		status string
		want   string
	}{
		{constants.ActivityStatusDraft, "草稿"},
		{constants.ActivityStatusPublished, "已发布"},
		{constants.ActivityStatusOngoing, "进行中"},
		{constants.ActivityStatusFinished, "已结束"},
		{constants.ActivityStatusCancelled, "已取消"},
		{"unknown", "未知"},
	}
	for _, tc := range cases {
		if got := ActivityStatusText(tc.status); got != tc.want {
			t.Fatalf("ActivityStatusText(%q) = %q, want %q", tc.status, got, tc.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		seconds int
		want    string
	}{
		{0, "--:--:--"},
		{3661, "01:01:01"},
		{125, "00:02:05"},
	}
	for _, tc := range cases {
		if got := FormatDuration(tc.seconds); got != tc.want {
			t.Fatalf("FormatDuration(%d) = %q, want %q", tc.seconds, got, tc.want)
		}
	}
}

func TestDifficultyText(t *testing.T) {
	cases := []struct {
		difficulty string
		want       string
	}{
		{constants.DifficultyFamily, "亲子"},
		{constants.DifficultyAdult, "成人"},
		{constants.DifficultyPro, "专业"},
	}
	for _, tc := range cases {
		if got := DifficultyText(tc.difficulty); got != tc.want {
			t.Fatalf("DifficultyText(%q) = %q, want %q", tc.difficulty, got, tc.want)
		}
	}
}

func TestMaskPhone(t *testing.T) {
	if got := MaskPhone("13812345678"); got != "138****5678" {
		t.Fatalf("MaskPhone = %q", got)
	}
}
