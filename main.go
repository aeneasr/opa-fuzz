package main

import (
	"os/exec"
	"io/ioutil"
	"fmt"
)

func main() {
	i := 0
	for l := 1; l <= 3; l++ {
		np := next(l)

		for {
			n := np()
			if i % 500 == 0 {
				fmt.Printf("Got next (%d) for %d at %d\n", len(n), l, i)
			}
			i++

			if len(n) == 0 {
				break
			}

			if err := ioutil.WriteFile("fuzz.rego", n, 0644); err != nil {
				panic(err)
			}

			c := exec.Command("opa", "run", "fuzz.rego")
			if err := c.Run(); err != nil {
				if err.Error() != "exit status 1" {
					fmt.Printf("Payload %s ( %x ) caused error %s\n", n, n, err)
				}
			}
		}
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
