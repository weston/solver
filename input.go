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

// Weights is a map from hand to a weight at which it appears in this range
type WeightedRange struct {
	Weights map[Hand]int
}

// Copy returns a copy of a WeightedRange
func (wr WeightedRange) Copy() WeightedRange {
	weights := make(map[Hand]int, len(wr.Weights))
	for k, v := range wr.Weights {
		weights[k] = v
	}
	return WeightedRange{Weights: weights}
}

type ActionType string

const (
	ActionTypeBet   ActionType = "bet"
	ActionTypeCall  ActionType = "call"
	ActionTypeCheck ActionType = "check"
	ActionTypeFold  ActionType = "fold"
	ActionTypeRaise ActionType = "raise"
)

type Action struct {
	Type   ActionType
	Amount int
}

// EffectiveBetAmount calculates the effective bet amount given the stack
// sizes and the bet percentage
func (a Action) EffectiveBetAmount(pot, effectiveStack int) (int, bool) {
	if a.Type == ActionTypeBet || a.Type == ActionTypeRaise {
		bet := (a.Amount * pot) / 100
		if effectiveStack < bet {
			return effectiveStack, true
		}
		return bet, false
	}
	return 0, false
}

func Bet(potPercentage int) Action {
	return Action{
		Type:   ActionTypeBet,
		Amount: potPercentage,
	}
}

func Check() Action {
	return Action{
		Type:   ActionTypeCheck,
		Amount: 0,
	}
}

func Raise(potPercentage int) Action {
	return Action{
		Type:   ActionTypeRaise,
		Amount: potPercentage,
	}
}

func Fold() Action {
	return Action{
		Type:   ActionTypeFold,
		Amount: 0,
	}
}

func Call() Action {
	return Action{
		Type:   ActionTypeCall,
		Amount: 0,
	}
}
