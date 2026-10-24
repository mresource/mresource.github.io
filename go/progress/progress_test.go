package progress

import (
	"math"
	"testing"
)

// Ani's example from the practice page: baseline 42, goal 70 after 12 weeks.
var ani = Plan{Baseline: 42, Goal: 70, Weeks: 12}

func TestAimAt(t *testing.T) {
	tests := []struct {
		name string
		plan Plan
		week float64
		want float64
	}{
		{"start", ani, 0, 42},
		{"end", ani, 12, 70},
		{"middle", ani, 6, 56},
		{"past the end keeps the slope", ani, 15, 77},
		{"flat goal", Plan{Baseline: 50, Goal: 50, Weeks: 8}, 4, 50},
		{"goal below baseline", Plan{Baseline: 60, Goal: 40, Weeks: 10}, 5, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plan.AimAt(tt.week); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("AimAt(%v) = %v, want %v", tt.week, got, tt.want)
			}
		})
	}
}

func TestDecide(t *testing.T) {
	tests := []struct {
		name      string
		plan      Plan
		scores    []float64
		changes   []int
		want      Decision
		wantSince int
	}{
		{"no points yet", ani, nil, nil, CollectMore, 0},
		{"three points are not enough", ani, []float64{40, 41, 42}, nil, CollectMore, 3},
		{"four below: change the approach", ani, []float64{43, 44, 45, 45}, nil, ChangeApproach, 4},
		{"four above: raise the goal", ani, []float64{50, 52, 55, 58}, nil, RaiseGoal, 4},
		{"mixed: continue", ani, []float64{43, 47, 45, 52}, nil, Continue, 4},
		// Aim at week 3 is exactly 49: a point on the line is neither below nor above.
		{"point on the line breaks a below run", ani, []float64{43, 44, 49, 45}, nil, Continue, 4},
		{"point on the line breaks an above run", ani, []float64{50, 52, 49, 58}, nil, Continue, 4},
		{"only the last four count", ani, []float64{30, 30, 30, 30, 60, 60, 60, 63}, nil, RaiseGoal, 8},
		{"counting restarts after a change", ani, []float64{43, 44, 45, 45, 52, 56}, []int{4}, CollectMore, 2},
		// Ani on the page: below for four weeks, change, then mostly above.
		{"page example after the change", ani, []float64{43, 44, 45, 45, 52, 56, 60, 63}, []int{4}, Continue, 4},
		{"latest change wins, order does not matter", ani, []float64{43, 44, 45, 45, 52, 56, 60, 63, 66}, []int{5, 2, 5}, Continue, 4},
		{"change at week zero is allowed", ani, []float64{43, 44, 45, 45}, []int{0}, ChangeApproach, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decide(tt.plan, tt.scores, tt.changes)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Decision != tt.want || got.SinceChange != tt.wantSince {
				t.Errorf("Decide() = %+v, want {%s %d}", got, tt.want, tt.wantSince)
			}
		})
	}
}

func TestDecideRejectsBadInput(t *testing.T) {
	tests := []struct {
		name    string
		plan    Plan
		scores  []float64
		changes []int
	}{
		{"zero weeks", Plan{Baseline: 40, Goal: 60, Weeks: 0}, nil, nil},
		{"negative baseline", Plan{Baseline: -1, Goal: 60, Weeks: 10}, nil, nil},
		{"NaN goal", Plan{Baseline: 40, Goal: math.NaN(), Weeks: 10}, nil, nil},
		{"negative score", ani, []float64{40, -3}, nil},
		{"infinite score", ani, []float64{math.Inf(1)}, nil},
		{"change after a week not measured", ani, []float64{40, 41}, []int{3}},
		{"negative change", ani, []float64{40, 41}, []int{-1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Decide(tt.plan, tt.scores, tt.changes); err == nil {
				t.Error("expected an error, got none")
			}
		})
	}
}
