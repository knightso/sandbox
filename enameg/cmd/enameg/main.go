package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/knightso/sandbox/enameg"
)

var (
	out = flag.String("out", "", "output file name; default {src}_name.go")
)

func main() {
	log.SetFlags(0)

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return
	}

	path := args[0]

	generated := enameg.Generate(path)
	if *out == "" {
		// 未指定の場合は {src}_name.go を出力先とする
		components := strings.Split(path, ".")
		*out = strings.Join(components[0:len(components)-1], ".") + "_name." + components[len(components)-1]
	}

	err := ioutil.WriteFile(*out, []byte(generated), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
