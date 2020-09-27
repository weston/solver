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

// CardFromString is a helper function to get a card object from a string
func CardFromString(s string) Card {
	//TODO: Validate some stuff
	return Card{
		Rank: string(s[0]),
		Suit: string(s[1]),
	}
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

// HandFromString is a helper function to get a hand object from a string
func HandFromString(s string) Hand {
	return Hand{
		c1: CardFromString(s[:2]),
		c2: CardFromString(s[2:]),
	}
}
