package main

import (
	"strconv"
	"strings"
)

// PreDefineSymbols 预定义符号
var PreDefineSymbols map[string]string

// dest=comp;jmp

// DestCodeTable 存储目标	代码表
var DestCodeTable map[string]string

// CompCodeTable 计算表达式 代码表
var CompCodeTable map[string]string

// JumpCodeTable 跳转状态	代码表
var JumpCodeTable map[string]string

// InitPreDefineSymbols 初始化预定义符号
func InitPreDefineSymbols() {
	PreDefineSymbols = make(map[string]string)
	// R0...R15
	for i := 0; i < 16; i++ {
		symbolValue := strconv.Itoa(i)
		symbolName := "R" + symbolValue
		PreDefineSymbols[symbolName] = symbolValue
	}
	// Screen,Keyboard
	PreDefineSymbols["SCREEN"] = "16384"
	PreDefineSymbols["KBD"] = "24576"
	// Other symbols
	otherSymbols := []string{"SP", "LCL", "ARG", "THIS", "THAT"}
	for i, symbolName := range otherSymbols {
		PreDefineSymbols[symbolName] = strconv.Itoa(i)
	}
}

// InitDestCodeTable 初始化存储目标代码表
func InitDestCodeTable() {
	DestCodeTable["NULL"] = "000"
	DestCodeTable["M"] = "001"
	DestCodeTable["D"] = "010"
	DestCodeTable["MD"] = "011"
	DestCodeTable["A"] = "100"
	DestCodeTable["AM"] = "101"
	DestCodeTable["AD"] = "110"
	DestCodeTable["AMD"] = "111"

}

// InitJumpCodeTable 初始化跳转状态代码表
func InitJumpCodeTable() {
	JumpCodeTable["NULL"] = "000"
	JumpCodeTable["JGT"] = "001"
	JumpCodeTable["JEQ"] = "010"
	JumpCodeTable["JGE"] = "011"
	JumpCodeTable["JLT"] = "100"
	JumpCodeTable["JNE"] = "101"
	JumpCodeTable["JLE"] = "110"
	JumpCodeTable["JMP"] = "111"
}

// InitCompCodeTable 初始化计算表达式代码表
func InitCompCodeTable() {
	CompCodeTable["0"] = "101010"
	CompCodeTable["1"] = "111111"
	CompCodeTable["-1"] = "111010"
	CompCodeTable["D"] = "001100"
	CompCodeTable["A"] = "110000"
	CompCodeTable["!D"] = "001101"
	CompCodeTable["!A"] = "110001"
	CompCodeTable["D+1"] = "011111"
	CompCodeTable["A+1"] = "110111"
	CompCodeTable["D-1"] = "001110"
	CompCodeTable["A-1"] = "110010"
	CompCodeTable["D+A"] = "000010"
	CompCodeTable["D-A"] = "010011"
	CompCodeTable["A-D"] = "000111"
	CompCodeTable["D&A"] = "000000"
	CompCodeTable["D|A"] = "010101"
	for comp, code := range CompCodeTable {
		// a=0,select A-Register
		CompCodeTable[comp] = "0" + code
		// a=1,select Memory input
		if strings.Contains(comp, "A") {
			MComp := strings.Replace(comp, "A", "M", -1)
			CompCodeTable[MComp] = "1" + code
		}
	}

}
