package object

import "fmt"

func argSizeCheck(expectedSize int, args []Object, body func([]Object) Object) Object {
	argsLength := len(args)
	if argsLength != expectedSize {
		return newError("wrong number of arguments. got=%d, want=%d", argsLength, expectedSize)
	}
	return body(args)
}

func arrayCheck(builtinName string, args []Object, body func(*Array, int) Object) Object {
	if args[0].Type() != ArrayObj {
		return newError("argument to '%s' must be ARRAY, got %s", builtinName, args[0].Type())
	}
	arr := args[0].(*Array)
	return body(arr, len(arr.Elements))
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			return argSizeCheck(1, args, func(args []Object) Object {
				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to 'len' not supported, got %s", arg.Type())
				}
			})
		}},
	},
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Printf(arg.Inspect())
			}

			return nil
		}},
	},
	{
		"first",
		&Builtin{Fn: func(args ...Object) Object {
			return argSizeCheck(1, args, func(objects []Object) Object {
				return arrayCheck("first", args, func(arr *Array, length int) Object {
					if length > 0 {
						return arr.Elements[0]
					}
					return nil
				})
			})
		}},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			return argSizeCheck(1, args, func(args []Object) Object {
				return arrayCheck("last", args, func(arr *Array, length int) Object {
					if length > 0 {
						return arr.Elements[length-1]
					}
					return nil
				})
			})
		}},
	},
	{
		"rest",
		&Builtin{Fn: func(args ...Object) Object {
			return argSizeCheck(1, args, func(args []Object) Object {
				return arrayCheck("rest", args, func(arr *Array, length int) Object {
					if length > 0 {
						newElements := make([]Object, length-1, length-1)
						copy(newElements, arr.Elements[1:length])
						return &Array{Elements: newElements}
					}
					return nil
				})
			})
		}},
	},
	{
		"push",
		&Builtin{Fn: func(args ...Object) Object {
			return argSizeCheck(2, args, func(args []Object) Object {
				return arrayCheck("push", args, func(arr *Array, length int) Object {
					newElements := make([]Object, length+1, length+1)
					copy(newElements, arr.Elements)
					newElements[length] = args[1]
					return &Array{Elements: newElements}
				})
			})
		}},
	},
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
