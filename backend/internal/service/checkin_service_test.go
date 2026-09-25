package service

import (
	"testing"
)

func TestHaversine(t *testing.T) {
	cases := []struct {
		name       string
		lat1, lng1 float64
		lat2, lng2 float64
		min, max   float64
	}{
		{"same point", 31.2304, 121.4737, 31.2304, 121.4737, 0, 1},
		{"approx 100m north", 31.2304, 121.4737, 31.2313, 121.4737, 80, 120},
		{"approx 1km", 31.2304, 121.4737, 31.2394, 121.4737, 900, 1200},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := haversine(tc.lat1, tc.lng1, tc.lat2, tc.lng2)
			if d < tc.min || d > tc.max {
				t.Fatalf("distance = %.1f m, want in [%.0f, %.0f]", d, tc.min, tc.max)
			}
		})
	}
}
