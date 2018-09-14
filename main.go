package main

import (
	"fmt"
	"flag"
	"github.com/open-policy-agent/opa/ast"
)

func main() {
	var chars = flag.Int("len", 4, "Maximum content length")
	var every = flag.Int("every", 100000, "Show status every X runs")
	flag.Parse()

	i := 0
	for l := 1; l < *chars; l++ {
		np := next(l)

		for {
			n := np()
			if i % *every == 0 {
				fmt.Printf("Got next (%d) for %d at %d\n", len(n), l, i)
			}
			i++

			if len(n) == 0 {
				break
			}

			run(n)
		}
	}
}

func run(module []byte) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Panic caused by \"%s\" was thrown with payload %x\n", r, module)
        }
    }()


	parsed, err := ast.ParseModule("example.rego", string(module))
	if err != nil {
		// This is fine
		// fmt.Printf("Unable to parse module with payload %s ( %x ): %s", module, module, err)
		return
	}

	compiler := ast.NewCompiler()
	compiler.Compile(map[string]*ast.Module{
	    "example.rego": parsed,
	})


	if compiler.Failed() {
		fmt.Printf("Unable to compile module with payload %s ( %x ): %s", module, module, compiler.Errors)
	}
}


func next(n int) func() []byte {
    p := make([]byte, n)
    x := make([]int, len(p))
    return func() []byte {
        p := p[:len(x)]
        for i, xi := range x {
            p[i] = byte(xi)
        }
        for i := len(x) - 1; i >= 0; i-- {
            x[i]++
            if x[i] < 255 {
                break
            }
            x[i] = 0
            if i <= 0 {
                x = x[0:0]
                break
            }
        }
        return p
    }
}
