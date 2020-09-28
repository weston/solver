package main

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

// CardFromString is a helper function to get a card object from a string
func CardFromString(s string) Card {
	//TODO: Validate some stuff
	return Card{
		Rank: string(s[0]),
		Suit: string(s[1]),
	}
}

// HandFromString is a helper function to get a hand object from a string
func HandFromString(s string) Hand {
	return Hand{
		c1: CardFromString(s[:2]),
		c2: CardFromString(s[2:]),
	}
}
