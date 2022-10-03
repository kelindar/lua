# Concurrent LUA Executor

[![Go Report Card](https://goreportcard.com/badge/github.com/kelindar/lua)](https://goreportcard.com/report/github.com/kelindar/lua)
[![GoDoc](https://godoc.org/github.com/kelindar/lua?status.svg)](https://godoc.org/github.com/kelindar/lua)

This repository contains a concurrent LUA executor that is designed to keep running a same (but updateable) set of scripts over a long period of time. The design of the library is quite opinionated, as it requires the `main()` function to be present in the script in order to run it. It also maintains a pool of VMs per `Script` instance to increase throughput per script.

Under the hood, it uses [gopher-lua](https://github.com/yuin/gopher-lua) library but abstracts it away, in order to be easily replaced in future if required. 


## Usage Example
Below is the usage example which runs the fibonacci LUA script with input `10`.

```go
// Load the script
s, err := FromString("test.lua", `
    function main(n)
        if n < 2 then return 1 end
        return main(n - 2) + main(n - 1)
    end
`)

// Run the main() function with 10 as argument
result, err := s.Run(context.Background(), 10)
println(result.String()) // Output: 89
```

The library also supports passing complex data types, thanks to [gopher-luar](https://github.com/layeh/gopher-luar). In the example below we create a `Person` struct and update its name in LUA as a side-effect of the script. It also returns the updated name back as a string.

```go
// Load the script
s, err := FromString("test.lua", `
    function main(input)
        input.Name = "Updated"
        return input.Name
    end
`)

input := &Person{ Name: "Roman" }
out, err := s.Run(context.Background(), input)
println(out)         // Outputs: "Updated"
println(input.Name)  // Outputs: "Updated"
```

## Native Modules

This library also supports and abstracts modules, which allows you to provide one or multiple native libraries which can be used by the script. These things are just ensembles of functions which are implemented in pure Go. 

Such functions must comply to a specific interface - they should have their arguments as the library's values (e.g. `Number`, `String` or `Bool`) and the result can be either a value and `error` or just an `error`. Here's an example of such function:
```go
func hash(s lua.String) (lua.Number, error) {
	h := fnv.New32a()
	h.Write([]byte(s))

	return lua.Number(h.Sum32()), nil
}
```

In order to use it, the functions should be registered into a `NativeModule` which then is loaded when script is created.
```go
// Create a test module which provides hash function
module := &NativeModule{
    Name:    "test",
    Version: "1.0.0",
}
module.Register("hash", hash)

// Load the script
s, err := FromString("test.lua", `
    local api = require("test")

    function main(input)
        return api.hash(input)
    end
`, module) // <- attach the module

out, err := s.Run(context.Background(), "abcdef")
println(out) // Output: 4282878506

```

## Script Modules 

Similarly to native modules, the library also supports LUA script modules. In order to use it, first you need to create a script which contains a module and returns a table with the functions. Then, create a `ScriptModule` which points to the script with `Name` which can be used in the `require` statement.

```go
moduleCode, err := FromString("module.lua", `
    local demo_mod = {} -- The main table

    function demo_mod.Mult(a, b)
        return a * b
    end

    return demo_mod
`)

// Create a test module which provides hash function
module := &ScriptModule{
    Script:  moduleCode,
    Name:    "demo_mod",
    Version: "1.0.0",
}
```

Finally, attach the module to the script as with native modules.

```go
// Load the script
s, err := FromString("test.lua", `
    local demo = require("demo_mod")

    function main(input)
        return demo.Mult(5, 5)
    end
`, module) // <- attach the module

out, err := s.Run(context.Background())
println(out) // Output: 25
```


## Benchmarks

```
Benchmark_Serial/fib-8         	 5870025	       203 ns/op	      16 B/op	       2 allocs/op
Benchmark_Serial/empty-8       	 8592448	       137 ns/op	       0 B/op	       0 allocs/op
Benchmark_Serial/update-8      	 1000000	      1069 ns/op	     224 B/op	      14 allocs/op
```
