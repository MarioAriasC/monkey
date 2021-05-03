package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	statements := func(exp Expression) []Statement{
		return []Statement{
			&ExpressionStatement{Expression: exp},
		}
	}

	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}
		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	test := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},
		{
			&Program{
				Statements: statements(one()),
			},
			&Program{
				Statements: statements(two()),
			},
		},
		{
			&InfixExpression{Left: one(), Operator: "+", Right: two()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&InfixExpression{Left: two(), Operator: "+", Right: one()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpression{Operator: "-", Right: one()},
			&PrefixExpression{Operator: "-", Right: two()},
		},
		{
			&IndexExpression{Left: one(), Index: one()},
			&IndexExpression{Left: two(), Index: two()},
		},
		{
			&IfExpression{
				Condition: one(),
				Consequence: &BlockStatement{
					Statements: statements(one()),
				},
				Alternative: &BlockStatement{
					Statements: statements(one()),
				},
			},
			&IfExpression{
				Condition: two(),
				Consequence: &BlockStatement{
					Statements: statements(two()),
				},
				Alternative: &BlockStatement{
					Statements: statements(two()),
				},
			},
		},
		{
			&ReturnStatement{ReturnValue: one()},
			&ReturnStatement{ReturnValue: two()},
		},
		{
			&LetStatement{Value: one()},
			&LetStatement{Value: two()},
		},
		{
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: statements(one()),
				},
			},
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: statements(two()),
				},
			},
		},
		{
			&ArrayLiteral{Elements: []Expression{one(), one()}},
			&ArrayLiteral{Elements: []Expression{two(), two()}},
		},
	}

	for _, tt := range test {
		modified := Modify(tt.input, turnOneIntoTwo)

		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expected)
		}
	}

	hashLiteral := &HashLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}

	Modify(hashLiteral, turnOneIntoTwo)

	for key, val := range hashLiteral.Pairs {
		key, _ := key.(*IntegerLiteral)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, key.Value)
		}

		val, _ := val.(*IntegerLiteral)
		if key.Value != 2 {
			t.Errorf("value is not %d, got=%d", 2, val.Value)
		}
	}
}
