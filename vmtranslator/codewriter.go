package main

import (
	"os"
)

var AsmSourceFile *os.File

func CreateAsmSourceFile(fileName string) (err error) {
	AsmSourceFile, err = os.Create(fileName)
	return
}

func CloseAsmSourceFile(fileName string) {
	AsmSourceFile.Close()
}

// WriteCommand write command to asm file
func WriteCommand(command *Command) {
	// STACK 	RAM:256-2047  	SP:R0=> 栈顶位置(RAM[0])
	// HEAP 	RAM:2048-16383
	// local    				LCL:local段基址(RAM[1])
	// static 	RAM:16-255
	// constant
	// argument					ARG:argument段基址(RAM[2])
	// pointer
	// this						THIS:RAM[3]
	// that						THAT:RAM[4]
	// temp						RAM[5..12]
	switch command.Type {
	case CEmpty:
		return
	case CArithmetic:
	case CPush:
	case CPop:
	}
}
