package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lshaowei18/codekata/args/go/parser"
)

func main() {
	p, err := parser.Args("l,p#,d*", os.Args)
	if err != nil {
		log.Fatalln(err)
	}

	v, _ := p.GetValue("d")
	fmt.Printf("l : %v", v)
}
