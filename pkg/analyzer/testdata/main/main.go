package main

import (
	"os"
)

func main() {
	os.Exit(0) // want "found os.Exit function call in main package, in main function"
}
