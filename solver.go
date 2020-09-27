package main

// StrategyNode is a node in the strategy tree of a river situation
type StrategyNode struct {
	Parent         *StrategyNode
	PreviousAction *Action
	Children       map[Action]StrategyNode
	Weights        map[Action]WeightedRange
	Pot            int
	LiveMoney      int
	OOPStack       int
	IPStack        int
}

func (sn StrategyNode) EffectiveStack() int {
	if sn.OOPStack < sn.IPStack {
		return sn.OOPStack
	}
	return sn.IPStack
}

// Builds a tree with no intelligent strategy suggestions
func BuildTree(input SolverInput) StrategyNode {
	sn := StrategyNode{
		Parent:         nil,
		PreviousAction: nil,
		Children:       make(map[Action]StrategyNode),
		Weights:        make(map[Action]WeightedRange),
		Pot:            input.Pot,
		OOPStack:       input.EffectiveStack,
		IPStack:        input.EffectiveStack,
		LiveMoney:      0,
	}
	for _, action := range input.OOPActions {
		if action.Type == ActionTypeRaise {
			continue
		}

		sn.Children[action] = buildTreeHelperIP(
			input, sn, action)
		sn.Weights[action] = input.OOPRange.Copy()
	}
	return sn
}

// parent will not have Weights and actions populated
func buildTreeHelperIP(input SolverInput, parent StrategyNode,
	previousAction Action) StrategyNode {

	sn := StrategyNode{
		Pot:            parent.Pot,
		LiveMoney:      0,
		OOPStack:       parent.OOPStack,
		IPStack:        parent.IPStack,
		Children:       make(map[Action]StrategyNode),
		Weights:        make(map[Action]WeightedRange),
		Parent:         &parent,
		PreviousAction: &previousAction,
	}
	isAllin := false
	var previousBetAmount int
	if previousAction.Type == ActionTypeBet {
		previousBetAmount, isAllin = previousAction.EffectiveBetAmount(
			parent.Pot, parent.EffectiveStack())
		sn.LiveMoney = previousBetAmount
		sn.OOPStack -= previousBetAmount
	} else if previousAction.Type == ActionTypeRaise {
		// Model raises as a call + another bet
		callAmount := parent.LiveMoney
		sn.Pot += 2 * callAmount
		sn.OOPStack -= callAmount
		previousBetAmount, isAllin = previousAction.EffectiveBetAmount(
			sn.Pot, sn.EffectiveStack())
		sn.LiveMoney = previousBetAmount
		sn.OOPStack -= previousBetAmount
	}

	validActionTypes := make(map[ActionType]bool)
	validActions := make([]Action, 0)
	if previousAction.Type == ActionTypeCheck {
		validActionTypes[ActionTypeCheck] = true
		validActionTypes[ActionTypeBet] = true
	} else {
		validActionTypes[ActionTypeRaise] = true
		validActions = append(validActions, Call(), Fold())
	}
	for _, action := range input.IPActions {
		if validActionTypes[action.Type] {
			validActions = append(validActions, action)
		}
	}
	for _, action := range validActions {
		if !isAllin && (action.Type == ActionTypeBet || action.Type == ActionTypeRaise) {
			sn.Children[action] = buildTreeHelperOOP(input, sn, action)
		}
		if action.Type == ActionTypeRaise && isAllin {
			// Don't set a strategic range for allin
			continue
		}
		sn.Weights[action] = input.IPRange.Copy()
	}
	return sn
}

// Only called when IP has bet or raised
// Valid actions are call, fold, and raise
func buildTreeHelperOOP(input SolverInput, parent StrategyNode,
	previousAction Action) StrategyNode {
	if previousAction.Type != ActionTypeBet &&
		previousAction.Type != ActionTypeRaise {

		panic("Invalid previous action type")
	}
	sn := StrategyNode{
		Pot:            parent.Pot,
		LiveMoney:      0,
		OOPStack:       parent.OOPStack,
		IPStack:        parent.IPStack,
		Children:       make(map[Action]StrategyNode),
		Weights:        make(map[Action]WeightedRange),
		Parent:         &parent,
		PreviousAction: &previousAction,
	}
	var betAmount int
	var isAllin bool
	if previousAction.Type == ActionTypeBet {
		betAmount, isAllin = previousAction.EffectiveBetAmount(parent.Pot,
			parent.EffectiveStack())
		sn.LiveMoney = betAmount
		sn.IPStack -= sn.LiveMoney
	} else if previousAction.Type == ActionTypeRaise {
		callAmount := parent.LiveMoney
		sn.Pot += 2 * callAmount
		sn.IPStack -= callAmount
		betAmount, isAllin = previousAction.EffectiveBetAmount(
			sn.Pot, sn.EffectiveStack())
		sn.LiveMoney = betAmount
		sn.IPStack -= betAmount
	}
	validActions := []Action{Call(), Fold()}
	for _, action := range input.OOPActions {
		if action.Type == ActionTypeRaise {
			validActions = append(validActions, action)
		}
	}
	for _, action := range validActions {
		if !isAllin && action.Type == ActionTypeRaise {
			sn.Children[action] = buildTreeHelperOOP(input, sn, action)
		}
		if action.Type == ActionTypeRaise && isAllin {
			// Don't set a strategic range for allin
			continue
		}
		sn.Weights[action] = input.OOPRange.Copy()
	}
	return sn
}
