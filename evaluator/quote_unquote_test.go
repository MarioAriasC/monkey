package evaluator

import (
	"monkey/object"
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(5)`,
			`5`,
		},
		{
			`quote(5 + 8)`,
			`(5 + 8)`,
		},
		{
			`quote(foobar)`,
			`foobar`,
		},
		{
			`quote(foobar + barfoo)`,
			`(foobar + barfoo)`,
		},
	}

	for _, tt := range tests {
		testQuote(t, tt.input, tt.expected)
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(unquote(4))`,
			`4`,
		},
		{
			`quote(unquote(4 + 4))`,
			`8`,
		},
		{
			`quote(8 + unquote(4 + 4))`,
			`(8 + 8)`,
		},
		{
			`quote(unquote(4 + 4) + 8)`,
			`(8 + 8)`,
		},
		{
			`let foobar = 8;
					quote(foobar)`,
			`foobar`,
		},
		{
			`let foobar = 8;
					quote(unquote(foobar))`,
			`8`,
		},
		{
			`quote(unquote(true))`,
			`true`,
		},
		{
			`quote(unquote(true == false))`,
			`false`,
		},
		{
			`quote(unquote(quote(4 + 4)))`,
			`(4 + 4)`,
		},
		{
			`let quotedInfixExpression = quote(4 + 4);
					quote(unquote(4 + 4) + unquote(quotedInfixExpression))`,
			`(8 + (4 + 4))`,
		},

	}

	for _, tt := range tests {
		testQuote(t, tt.input, tt.expected)
	}
}

func testQuote(t *testing.T, input, expected string) {
	evaluated := testEval(input)
	quote, ok := evaluated.(*object.Quote)
	if !ok {
		t.Fatalf("expected *object.Quote. got=%T (%+v)", evaluated, evaluated)
	}

	if quote.Node == nil {
		t.Fatalf("quote.Node is nil")
	}

	if quote.Node.String() != expected {
		t.Errorf("not equal. got=%q, want=%q", quote.Node.String(), expected)
	}
}
