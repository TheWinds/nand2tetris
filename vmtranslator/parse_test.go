package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommandParse(t *testing.T) {
	deepEqual := func(commandStr string, b interface{}) bool {
		return reflect.DeepEqual(*new(Command).parse(commandStr), b)
	}
	Convey("test command parser", t, func() {
		arithmetics := []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}
		for _, commandName := range arithmetics {
			So(deepEqual(commandName,
				Command{Name: commandName, FirstArg: commandName, Type: CArithmetic}),
				ShouldBeTrue)
		}
		So(deepEqual("push constant 1",
			Command{Name: "push", FirstArg: "constant", SecondArg: 1, Type: CPush}),
			ShouldBeTrue)
		So(deepEqual("pop constant 1",
			Command{Name: "pop", FirstArg: "constant", SecondArg: 1, Type: CPop}),
			ShouldBeTrue)
	})
}

func TestFileParse(t *testing.T) {
	file, _ := ioutil.TempFile("", "test_vm_file")
	vmCode := `// This file is part of www.nand2tetris.org
	// and the book "The Elements of Computing Systems"
	// by Nisan and Schocken, MIT Press.
	// File name: projects/07/MemoryAccess/BasicTest/BasicTest.vm
	
	// Executes pop and push commands using the virtual memory segments.
	push constant 10
	pop local 0
	push constant 21
	push constant 22
	pop argument 2
	pop argument 1
	push constant 36
	pop this 6
	push constant 42
	push constant 45
	pop that 5
	pop that 2
	push constant 510
	pop temp 6
	push local 0
	push that 5
	add
	push argument 1
	sub
	push this 6
	push this 6
	add
	sub
	push temp 6
	add`
	fileName := file.Name()
	fmt.Println(fileName)
	file.WriteString(vmCode)
	file.Close()
	sc, _ := OpenVMSourcrFile(fileName)
	for HasMoreCommands(sc) {
		Advance(sc)
		fmt.Println(*CurrentCommand)
	}
	CloseVMSourcrFile()
	os.Remove(fileName)
}

func Test1(t *testing.T) {
	// fmt.Printf(genCmp("gt"))
	// fmt.Printf(genCmp("lt"))
	// fmt.Printf(genCmp("eq"))
	fmt.Printf(genComp("not"))
}
