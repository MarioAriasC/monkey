#Monkey 

A go implementation of the [Monkey Language](https://monkeylang.org/)

Following both books, it implements the interpreter and the compiler

My code is shorter as I used high order functions in some places to reduce code duplication, i.e.:

## Original

```go
var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		default:
			return newError("argument to `len` not supported, got %s",
				args[0].Type())
		}
	},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
```

## Refactored

```go
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
```

# Running

Build the executable with:

```shell
$ go build -o monkey .
```

And run it with:

```shell
$ ./monkey
```

# Benchmarks

Build the benchmark with:

```shell
$ go build -o fibonacci ./benchmark
```

And run it with:

```shell
$ ./fibonacci
```

You can pass two additional parameters:

 - `engine` with two possible values: `vm` and `eval`.
 - `algo` with two possible values: `slow` and `fast`.

# Tests

Run the tests with:

```shell
$ go test ./...
```