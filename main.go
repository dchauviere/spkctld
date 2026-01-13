package main

import (
	"fmt"
	"os"

	"github.com/dchauviere/spkctld/cmd"
)

func main() {
	if err := (&cmd.RootCmd{}).Command().Execute(); err != nil {
		fmt.Printf("fail to execute : %s\n", err)
		os.Exit(1)
	}
}
