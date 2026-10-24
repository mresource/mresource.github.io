// Package progress implements the progress-monitoring logic behind mresource:
// an aim line from a baseline to a goal, and the four-point decision rule used
// in data-based individualization (DBI).
//
// The rule looks only at measurements taken after the most recent change of
// approach. When the last four of them are all below the aim line, it is time
// to change the approach. When all four are above, the goal can be raised.
// A point exactly on the aim line counts as neither below nor above.
package progress

import (
	"errors"
	"fmt"
	"math"
)

// PointsForDecision is how many consecutive measurements the rule needs.
const PointsForDecision = 4

// Plan describes where a child starts and where the team aims to be.
type Plan struct {
	Baseline float64 `json:"baseline"` // starting score, e.g. words correct per minute
	Goal     float64 `json:"goal"`     // target score at the end of the plan
	Weeks    int     `json:"weeks"`    // weeks from baseline to goal
}

// Validate reports whether the plan can draw an aim line.
func (p Plan) Validate() error {
	switch {
	case p.Weeks <= 0:
		return errors.New("weeks must be greater than zero")
	case !finiteNonNeg(p.Baseline):
		return errors.New("baseline must be a number zero or above")
	case !finiteNonNeg(p.Goal):
		return errors.New("goal must be a number zero or above")
	}
	return nil
}

// AimAt returns the expected score at the given week on the straight line
// from Baseline (week 0) to Goal (week Weeks). Past Weeks the line continues
// with the same slope.
func (p Plan) AimAt(week float64) float64 {
	return p.Baseline + (p.Goal-p.Baseline)*week/float64(p.Weeks)
}

// Decision is what the rule suggests after the latest measurement.
type Decision string

const (
	CollectMore    Decision = "collect_more"    // fewer than four points since the last change
	Continue       Decision = "continue"        // points sit around the aim line
	ChangeApproach Decision = "change_approach" // four in a row below the aim line
	RaiseGoal      Decision = "raise_goal"      // four in a row above the aim line
)

// Result is the outcome of Decide.
type Result struct {
	Decision    Decision `json:"decision"`
	SinceChange int      `json:"sinceChange"` // measurements since the last change
}

// Decide applies the four-point rule.
//
// scores[i] is the measurement of week i+1. changes lists the weeks after
// which the approach was changed (a change after week 4 means weeks 5, 6, …
// belong to the new approach). Order and duplicates in changes do not matter.
func Decide(p Plan, scores []float64, changes []int) (Result, error) {
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	for i, v := range scores {
		if !finiteNonNeg(v) {
			return Result{}, fmt.Errorf("week %d: score must be a number zero or above", i+1)
		}
	}
	last := 0
	for _, c := range changes {
		if c < 0 || c > len(scores) {
			return Result{}, fmt.Errorf("change after week %d is outside the %d measured weeks", c, len(scores))
		}
		if c > last {
			last = c
		}
	}

	since := len(scores) - last
	res := Result{SinceChange: since}
	if since < PointsForDecision {
		res.Decision = CollectMore
		return res, nil
	}

	below, above := 0, 0
	for week := len(scores) - PointsForDecision + 1; week <= len(scores); week++ {
		v, aim := scores[week-1], p.AimAt(float64(week))
		switch {
		case v < aim:
			below++
		case v > aim:
			above++
		}
	}
	switch {
	case below == PointsForDecision:
		res.Decision = ChangeApproach
	case above == PointsForDecision:
		res.Decision = RaiseGoal
	default:
		res.Decision = Continue
	}
	return res, nil
}

func finiteNonNeg(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0
}
