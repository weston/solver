package main

import (
	"testing"
)

func TestTreeCheckOnly(t *testing.T) {
	board, oopRange, ipRange := setUp()
	oopActions := []Action{Check()}
	ipActions := []Action{Check()}

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
	root, err := BuildTree(input)
	if err != nil {
		t.Error(err)
	}
	numNodes := treeSize(root)
	if treeSize(root) != 2 {
		t.Errorf("Expected 2 nodes, got %v", numNodes)
	}
	validateNodeChipCounts(t, root, 300, 0, 700, 700)
}

func TestTreeIPBetOrCheckNoRaise(t *testing.T) {
	board, oopRange, ipRange := setUp()
	oopActions := []Action{Check()}
	ipActions := []Action{Bet(10), Check()}
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
	root, err := BuildTree(input)
	if err != nil {
		t.Error(err)
	}
	numNodes := treeSize(root)
	if treeSize(root) != 3 {
		t.Errorf("Expected 3 nodes, got %v", numNodes)
	}
	if len(root.Children) != 1 {
		t.Errorf("Expected 1 child, got %v", len(root.Children))
	}
	validateNodeChipCounts(t, root, 300, 0, 700, 700)

	ipNode, exists := root.Children[Check()]
	if !exists {
		t.Errorf("Did not find a Check node in OOP Root")
	}
	validateNodeChipCounts(t, ipNode, 300, 0, 700, 700)

	if len(ipNode.Children) != 1 {
		t.Errorf("Expected 1 child, got %v", len(ipNode.Children))
	}

	oopNode, exists := ipNode.Children[Bet(10)]
	if !exists {
		t.Error("Call/Fold node for OOP does not exist")
	}
	if len(oopNode.Children) != 0 {
		t.Errorf("OOP call/fold node had %v unexpected children",
			len(oopNode.Children))
	}
	validateNodeChipCounts(t, oopNode, 300, 30, 700, 670)
}

func TestIPBetAllin(t *testing.T) {
	board, oopRange, ipRange := setUp()
	oopActions := []Action{Check()}
	ipActions := []Action{Bet(100), Check()}
	input := SolverInput{
		Board:          board,
		OOPActions:     oopActions,
		IPActions:      ipActions,
		OOPRange:       oopRange,
		IPRange:        ipRange,
		Pot:            100,
		EffectiveStack: 100,
		JamThreshold:   90,
	}
	root, err := BuildTree(input)
	if err != nil {
		t.Error(err)
	}
	numNodes := treeSize(root)
	if treeSize(root) != 3 {
		t.Errorf("Expected 3 nodes, got %v", numNodes)
	}
	validateNodeChipCounts(t, root, 100, 0, 100, 100)
	ipNode, exists := root.Children[Check()]
	if !exists {
		t.Error("Missing node for IP")
	}
	validateNodeChipCounts(t, ipNode, 100, 0, 100, 100)
	oopNode, exists := ipNode.Children[Bet(100)]
	if !exists {
		t.Error("Missing node for OOP")
	}
	validateNodeChipCounts(t, oopNode, 100, 100, 100, 0)
}

func setUp() ([]Card, WeightedRange, WeightedRange) {
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
		},
	}
	ipRange := WeightedRange{
		Weights: map[Hand]int{
			HandFromString("AhAd"): 100,
		},
	}
	return board, oopRange, ipRange
}

func treeSize(root StrategyNode) int {
	childrenSum := 0
	for _, child := range root.Children {
		childrenSum += treeSize(child)
	}
	return childrenSum + 1
}

func validateNodeChipCounts(t *testing.T, node StrategyNode, pot, liveMoney,
	oopStack, ipStack int) {

	if node.Pot != pot {
		t.Errorf("expected pot of %v got %v", pot, node.Pot)
	}
	if node.LiveMoney != liveMoney {
		t.Errorf("Expected LiveMoney of %v got %v", liveMoney, node.LiveMoney)
	}
	if node.IPStack != ipStack {
		t.Errorf("Expected IPStack of %v got %v", ipStack, node.IPStack)
	}
	if node.OOPStack != oopStack {
		t.Errorf("Expected OOPStack of %v got %v", oopStack, node.OOPStack)
	}

}
