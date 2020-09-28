package main

import (
	"errors"
)

// SolverInput is the input to the river solver
type SolverInput struct {
	OOPRange       WeightedRange
	OOPActions     []Action
	IPRange        WeightedRange
	IPActions      []Action
	Pot            int
	EffectiveStack int
	Board          []Card
	JamThreshold   int
}

// Validate validates the solver input
func (s SolverInput) Validate() error {
	//TODO: Validate more
	if len(s.Board) != 5 {
		return errors.New("Invalid number of cards on the board")
	}
	return nil
}
