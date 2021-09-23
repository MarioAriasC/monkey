package main

import (
	"flag"
	"fmt"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"time"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
var algo = flag.String("algo", "fast", "use 'slow' or 'fast'")

var slowInput = `
let fibonacci = fn(x) {
	if (x == 0) {
		return 0;
	} else {
		if (x == 1) {
			return 1;
		} else {
			fibonacci(x - 1) + fibonacci(x - 2);
		}
	}
};
fibonacci(35);
`

var fastInput = `
let fibonacci = fn(x) {    	
    if (x < 2) {
    	return x;
    } else {
    	fibonacci(x - 1) + fibonacci(x - 2);
    }
};
fibonacci(35);        
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	var input string
	if *algo == "slow" {
		input = slowInput
	} else {
		input = fastInput
	}

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if *engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}
		machine := vm.New(comp.Bytecode())

		start := time.Now()
		err = machine.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf("engine=%s, algo=%s, result=%s, duration=%s\n", *engine, *algo, result.Inspect(), duration)
}
