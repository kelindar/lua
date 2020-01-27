# Concurrent LUA Executor

This repository contains a concurrent LUA executor that is designed to keep running a same (but updateable) set of scripts over a long period of time. The design of the library is quite opinionated, as it requires the `main()` function to be present in the script in order to run it. It also maintains a single VM per `Script` instance, protected by a mutex.

Under the hood, it uses [gopher-lua](https://github.com/yuin/gopher-lua) library but abstracts it away, in order to be easily replaced in future if required. 


## Usage Example
Below is the usage example which runs the fibonacci LUA script with input `10`.

```
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

```
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


## Benchmarks

```
Benchmark_Serial/fib-8         	 5870025	       203 ns/op	      16 B/op	       2 allocs/op
Benchmark_Serial/empty-8       	 8592448	       137 ns/op	       0 B/op	       0 allocs/op
Benchmark_Serial/update-8      	 1000000	      1069 ns/op	     224 B/op	      14 allocs/op
```