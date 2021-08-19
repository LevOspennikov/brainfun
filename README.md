# Brainfun

This is brainfuck language interpreter library. There is no look ahead in the implementation, it memorises input instead to return later.
Stack is used for functions, can be used basically for anything. The interpreter is very extensible, because functions can change state as they want.
However, it can not add symbols to input of the `Execute` and can not change history. State is also protected from Interpreter itself.

## Usage
```go
brainfuck := brainfun.NewInterpreter() // creates new interpreter with default set of function
```


```go
err := brainfuck.Execute(strings.NewReader("+++")) // executes program, should be ended, for example no open loops
```

You can replace functions of read and write with new, which contains your input 
```go 
buf := &bytes.Buffer{}
brainfuck.AddFunc('.', PrintWith(buf))
brainfuck.AddFunc(',', ReadWith(strings.NewReader("DBAC")))
```

Last - you can remove your func
```go
brainfuck.DeleteFunc('.')
```

## Possible improvements
* More tests 
* State can be encapsulated with Public API even for functions, but this would make extensions much more complex (but still possible, because Brainfuck is turing complete)
