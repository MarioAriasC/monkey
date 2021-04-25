package evaluator

import (
	"fmt"
	"monkey/object"
)

func argSizeCheck(expectedSize int, args []object.Object, body func([]object.Object) object.Object) object.Object {
	argsLength := len(args)
	if argsLength != expectedSize {
		return newError("wrong number of arguments. got=%d, want=1", argsLength)
	}

	return body(args)
}

func arrayCheck(builtinName string, args []object.Object, body func(*object.Array, int) object.Object) object.Object {
	if args[0].Type() != object.ArrayObj {
		return newError("argument to `%s` must be ARRAY, got %s", builtinName, args[0].Type())
	}
	arr := args[0].(*object.Array)
	return body(arr, len(arr.Elements))
}

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {

			return argSizeCheck(1, args, func(args []object.Object) object.Object {
				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", arg.Type())
				}

			})
		}},
	"first": {
		Fn: func(args ...object.Object) object.Object {

			return argSizeCheck(1, args, func(args []object.Object) object.Object {
				return arrayCheck("first", args, func(arr *object.Array, length int) object.Object {
					if length > 0 {
						return arr.Elements[0]
					}
					return NULL
				})

			})
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			return argSizeCheck(1, args, func(args []object.Object) object.Object {
				return arrayCheck("last", args, func(arr *object.Array, length int) object.Object {
					if length > 0 {
						return arr.Elements[length-1]
					}
					return NULL
				})
			})
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			return argSizeCheck(1, args, func(args []object.Object) object.Object {
				return arrayCheck("rest", args, func(arr *object.Array, length int) object.Object {
					if length > 0 {
						newElements := make([]object.Object, length-1, length-1)
						copy(newElements, arr.Elements[1:length])
						return &object.Array{Elements: newElements}
					}

					return NULL
				})
			})
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			return argSizeCheck(2, args, func(args []object.Object) object.Object {
				return arrayCheck("push", args, func(arr *object.Array, length int) object.Object {
					newElements := make([]object.Object, length+1, length+1)
					copy(newElements, arr.Elements)
					newElements[length] = args[1]

					return &object.Array{Elements: newElements}
				})
			})
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
