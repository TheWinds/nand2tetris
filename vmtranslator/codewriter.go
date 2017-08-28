package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var AsmSourceFile *os.File
var FileName string

func CreateAsmSourceFile(fileName string) (err error) {
	AsmSourceFile, err = os.Create(fileName)
	FileName = filepath.Base(fileName)
	FileName = strings.Replace(FileName, ".asm", "", -1)
	return
}

func CloseAsmSourceFile() {
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
	var asmCode string
	switch command.Type {
	case CEmpty:
		return
	case CArithmetic:
		switch command.Name {
		case "gt", "lt", "eq":
			asmCode = genCmp(command.Name)
		default:
			asmCode = genComp(command.Name)
		}
	case CPush:
		asmCode = genPush(command.FirstArg, command.SecondArg)
	case CPop:
		asmCode = genPop(command.FirstArg, command.SecondArg)
	}
	asmCode += "\n"
	// asmCode = "//" + command.Name + " " + command.FirstArg + " " + strconv.Itoa(command.SecondArg) + "\n" + asmCode
	AsmSourceFile.WriteString(asmCode)
}

// generate asm code that increase stack pointer
func genIncSP() string {
	return "@SP\nM=M+1"
}

// generate asm code that decrease stack pointer
// and pop value to D-Register
func genDecSP() string {
	return "@SP\nAM=M-1"
}

var cmpLabelCounter int64

// generate asm code that compare two elemnet
func genCmp(op string) string {
	cmpLabelCounter++
	endCmpLabel := fmt.Sprintf("END_CMP_%d", cmpLabelCounter)
	jump := ""
	switch op {
	case "eq":
		jump = "JEQ"
	case "lt":
		jump = "JLT"
	case "gt":
		jump = "JGT"
	}
	code := joinCode(
		genDecSP(),
		"D=M",
		genDecSP(),
		"D=M-D",
		"M=-1",
		"@"+endCmpLabel,
		"D;"+jump,
		"@SP",
		"A=M",
		"M=0",
		"("+endCmpLabel+")",
		genIncSP(),
	)
	return code
}

// generate asm code that compute two elemnet
func genComp(op string) string {
	var header string
	var midCode string
	if op != "neg" && op != "not" {
		header = joinCode(
			genDecSP(),
			"D=M",
			genDecSP(),
		)
	} else {
		header = genDecSP()
	}
	switch op {
	case "add":
		midCode = "M=M+D"
	case "sub":
		midCode = "M=M-D"
	case "and":
		midCode = "M=M&D"
	case "or":
		midCode = "M=M|D"
	case "neg":
		midCode = "M=-M"
	case "not":
		midCode = "M=!M"
	}
	return joinCode(
		header,
		midCode,
		genIncSP(),
	)

}

var segmentSymbolMap = map[string]string{
	"temp":     "5",
	"pointer":  "3",
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
}

// generate asm code that push (sengment index) to stack
func genPush(segment string, index int) string {
	var prepareCode string
	// locate address and put value to D-Register
	switch segment {
	case "constant":
		prepareCode = joinCode(
			fmt.Sprintf("@%d", index),
			"D=A",
		)
	case "static":
		prepareCode = joinCode(
			fmt.Sprintf("@%s.%d", FileName, index),
			"D=M",
		)
	case "temp", "pointer":
		baseAddress := 3
		if segment == "temp" {
			baseAddress = 5
		}
		prepareCode = joinCode(
			fmt.Sprintf("@%d", baseAddress+index),
			"D=M",
		)
	default:
		prepareCode = joinCode(
			fmt.Sprintf("@%d", index),
			"D=A",
			fmt.Sprintf("@%s", segmentSymbolMap[segment]),
			"A=D+M",
			"D=M",
		)
	}
	return joinCode(
		prepareCode,
		// push to stack
		"@SP",
		"A=M",
		"M=D",
		genIncSP(),
	)
}

// generate asm code that pop to (sengment index)
func genPop(segment string, index int) string {
	var prepareCode string
	// locate address and put address to R13
	switch segment {
	// if segment == (temp or pointer)
	// locate pointer address
	// else
	// locate pointer-to address
	case "temp", "pointer":
		prepareCode = joinCode(
			fmt.Sprintf("@%d", index),
			"D=A",
			fmt.Sprintf("@%s", segmentSymbolMap[segment]),
			"D=D+A",
			"@R13",
			"M=D")
	case "static":
		prepareCode = joinCode(
			fmt.Sprintf("@%s.%d", FileName, index),
			"D=A",
			"@R13",
			"M=D")
	default:
		prepareCode = joinCode(
			fmt.Sprintf("@%d", index),
			"D=A",
			fmt.Sprintf("@%s", segmentSymbolMap[segment]),
			"A=M",
			"D=D+A",
			"@R13",
			"M=D")
	}
	return joinCode(
		prepareCode,
		// pop to (segment index)
		genDecSP(),
		"D=M",
		"@R13",
		"A=M",
		"M=D")
}

func joinCode(code ...string) string {
	return strings.Join(code, "\n")
}
