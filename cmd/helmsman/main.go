package main

import (
	"os"

	"github.com/mkubaczyk/helmsman/internal/app"
)

func main() {
	exitCode := app.Main()
	os.Exit(exitCode)
}
