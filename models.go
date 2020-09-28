package main

import "errors"

var cardRanks = map[string]bool{
	"A": true,
	"2": true,
	"3": true,
	"4": true,
	"5": true,
	"6": true,
	"7": true,
	"8": true,
	"9": true,
	"T": true,
	"J": true,
	"Q": true,
	"K": true,
}

var cardSuits = map[string]bool{
	"c": true,
	"d": true,
	"h": true,
	"s": true,
}

// Card is a playing card
type Card struct {
	Rank string
	Suit string
}

func (c Card) Validate() error {
	if !cardRanks[c.Rank] {
		return errors.New("Invalid rank")
	}
	if !cardSuits[c.Suit] {
		return errors.New("Invalid suit")
	}
	return nil
}

func (c Card) Equals(other Card) bool {
	return c.Rank == other.Rank && c.Suit == other.Suit
}

// Hand is a 2 card poker hand
type Hand struct {
	c1 Card
	c2 Card
}

func (h Hand) Validate() error {
	err := h.c1.Validate()
	if err != nil {
		return err
	}
	err = h.c2.Validate()
	if err != nil {
		return err
	}
	if h.c1.Equals(h.c2) {
		return errors.New("Same card twice")
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
