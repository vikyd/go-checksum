package main

import (
	"encoding/json"
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
		fmt.Println("directory: " + file)
		doDir(file)
	case mode.IsRegular():
		fmt.Println("file: " + file)
		doGoMod(file)
	}
}

func doDir(dir string) {
	if len(os.Args) < 3 {
		fmt.Println("needs parameters: module path with version like: github.com/gin-gonic/gin@v1.4.0")
		return
	}

	prefix := os.Args[2]

	h, err := checksum.HashDir(dir, prefix)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(PrettyPrint(h))

}

func doGoMod(file string) {
	h, err := checksum.HashGoMod(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(PrettyPrint(h))

}

// PrettyPrint convert struct to pretty string
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
