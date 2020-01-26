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

## Future improvements

Future improvements might include:
* Maintaining a pool of VMs per Script instance for highly concurrent script execution.
* Support for input structs and tables
* Support for output structs and tables


## Benchmarks

```
Benchmark_Fib_Serial-8   	 4557145	       257 ns/op	      40 B/op	       3 allocs/op
```