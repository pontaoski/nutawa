// +build !wasm

package main

import (
	"Tawa/codegenerierung"
	"Tawa/interpreter"
	"Tawa/parser"
	"Tawa/typisierung"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alecthomas/repr"
)

func main() {
	if os.Args[1] == "--parser" {
		println(parser.Parser.String())
		return
	}
	datei, feh := ioutil.ReadFile(os.Args[1])
	if feh != nil {
		panic(feh)
	}
	es := parser.Datei{}
	feh = parser.Parser.ParseBytes(os.Args[1], datei, &es)
	if feh != nil {
		panic(feh)
	}
	es.Vorverarbeiten()
	v := typisierung.NeuVollKontext()
	v.Push()
	err := typisierung.Typisierung(v, &es)
	if err != nil {
		if v, ok := err.(*typisierung.Fehler); ok {
			s := string(datei)
			fmt.Printf("%s:%d:\n", os.Args[1], v.Line)
			println(strings.Split(s, "\n")[v.Line-1])
		}
		println(err.Error())
		os.Exit(1)
	}

	println(codegenerierung.Codegen(&es))

	vk := interpreter.NeuVollKontext()
	vk.Push()

	repr.Println(interpreter.Interpret(es, "main", vk))
}
