package main

import (
	"os"

	"sysconf-parser/cmd"
)

func main() {
	if len(os.Args) < 3 {
		cmd.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	var err error

	switch command {
	case "decode":
		err = cmd.Decode(filename)
	case "encode":
		err = cmd.Encode(filename)
	default:
		cmd.Usage()
		os.Exit(1)
	}

	if err != nil {
		panic(err)
	}
}
