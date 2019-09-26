package main

import (
	"fmt"

	"os"

	"github.com/vikyd/go-checksum/checksum"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("needs parameters: go.mod or dir")
		return
	}

	file := os.Args[1]
	fi, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		fmt.Println("directory")
		return
	case mode.IsRegular():
		// do file stuff
		fmt.Println("file")
	}

	// h, err := checksum.HashGoMod("/Users/viky/tmp/gosum/checksum-dir/gin/go.mod")
	h, err := checksum.HashGoMod("x/go.mod")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(checksum.PrettyPrint(h))

}
