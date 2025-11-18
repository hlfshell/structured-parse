// +build js wasm

package main

import (
	"github.com/hlfshell/structured-parse/go/structuredparse"
)

func main() {
	// Register all WASM exported functions
	structuredparse.RegisterWasmFunctions()
	
	// Keep the program running
	select {}
}


