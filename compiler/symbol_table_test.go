package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	testSymbol(t, "a", global, expected)
	testSymbol(t, "b", global, expected)
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range expected{
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func testSymbol(t *testing.T,name string, global *SymbolTable, expected map[string]Symbol) {
	symbol := global.Define(name)
	expectedSymbol := expected[name]
	if symbol != expectedSymbol {
		t.Errorf("expected a=%+v, got=%+v", expectedSymbol, symbol)
	}
}


