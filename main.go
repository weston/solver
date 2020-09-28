package main

import "fmt"

func main() {
	fmt.Println("Starting solver")
	board := []Card{
		CardFromString("Ac"),
		CardFromString("Kd"),
		CardFromString("9h"),
		CardFromString("2c"),
		CardFromString("2d"),
	}
	oopRange := WeightedRange{
		Weights: map[Hand]int{
			HandFromString("QhQc"): 100,
			HandFromString("QhQd"): 100,
			HandFromString("QhQs"): 100,
			HandFromString("QsQd"): 100,
			HandFromString("AhJh"): 100,
			HandFromString("AdJd"): 100,
			HandFromString("AhTh"): 100,
			HandFromString("AdTd"): 100,
		},
	}
	ipRange := WeightedRange{
		Weights: map[Hand]int{
			HandFromString("AhAd"): 100,
			HandFromString("KhKc"): 100,
			HandFromString("AhQh"): 100,
			HandFromString("AdQd"): 100,
			HandFromString("Ah2h"): 100,
			HandFromString("JhJc"): 100,
			HandFromString("JdJc"): 100,
			HandFromString("JsJc"): 100,
			HandFromString("JsJd"): 100,
		},
	}
	oopActions := []Action{
		Raise(50),
		Bet(50),
		Check(),
	}

	ipActions := []Action{
		Check(),
		Bet(25),
		Bet(75),
		Raise(50),
	}

	input := SolverInput{
		Board:          board,
		OOPActions:     oopActions,
		IPActions:      ipActions,
		OOPRange:       oopRange,
		IPRange:        ipRange,
		Pot:            300,
		EffectiveStack: 700,
		JamThreshold:   90,
	}
	root, _ := BuildTree(input)
	visualizeTree(root, "")
}

func visualizeTree(node StrategyNode, prefix string) {
	fmt.Println(prefix, "Pot:", node.Pot, "+", node.LiveMoney, "=", node.Pot+node.LiveMoney)
	fmt.Println(prefix, "OOP:", node.OOPStack, "   IP:", node.IPStack)
	fmt.Println(prefix, "Total:", node.OOPStack+node.IPStack+node.LiveMoney+node.Pot)
	for action, _ := range node.Weights {
		childNode, hasChildren := node.Children[action]
		if !hasChildren {
			fmt.Println(prefix, ">>>", action)
			continue
		}
		pot := node.Pot
		if action.Type == ActionTypeRaise {
			prevBet, _ := node.PreviousAction.EffectiveBetAmount(
				node.Parent.Pot+2*node.LiveMoney, node.Parent.EffectiveStack())
			pot += prevBet
		}
		rawBet, _ := action.EffectiveBetAmount(pot, node.EffectiveStack())
		fmt.Println(prefix, ">>>", action, rawBet)
		visualizeTree(childNode, prefix+"\t")
	}
}
