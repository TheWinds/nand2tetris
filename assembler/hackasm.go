package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if !checkArgs(args) {
		return
	}
	fileName := args[0]
	fileBytes, err := readFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	ts := time.Now()
	codes, err := ParseFile(fileBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	ReplaceSymbols(codes)
	if err != nil {
		fmt.Println(err)
		return
	}
	binaryCodes, err := Parser2Binary(codes)
	if err != nil {
		fmt.Println(err)
		return
	}
	hackFileName := getHackFileName(fileName)
	err = outputFile(hackFileName, binaryCodes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("compile success %s [%dms] \n", hackFileName, time.Since(ts).Nanoseconds()/1000000)
}

// read file
func readFile(fileName string) ([]byte, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("can not read file :" + err.Error())
	}
	return file, nil
}

func outputFile(fileName, content string) error {
	// if err := os.Stat(fileName); err != nil {
	// 	os.
	// }
	return ioutil.WriteFile(fileName, []byte(content), 0664)

}

func checkArgs(args []string) bool {
	if len(args) == 0 {
		fmt.Println("please use command 'hackasm [filename]' to compile file")
		return false
	}
	fileName := args[0]
	if path.Ext(fileName) != ".asm" {
		fmt.Println("only support .asm file")
		return false
	}
	return true
}

func getHackFileName(fileName string) string {
	return strings.Replace(fileName, ".asm", ".hack", -1)
}
