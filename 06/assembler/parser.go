package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CodeType type of code
type CodeType int

const (
	_ CodeType = iota
	CodeTypeUnknow
	CodeTypeInsA
	CodeTypeInsC
	CodeTypeJumpLabel
	CodeTypeCommentOrEmpty
)

// Code line of code
type Code struct {
	LineNum       int      //line number
	SourceLineNum int      //line number in source file
	Source        string   //source code
	Type          CodeType //type: A-Instruction,C-Instruction,Comment,Empty
}

func getJumpLable(code string) string {
	code = strings.Replace(code, " ", "", -1)
	return code[1 : len(code)-1]
}

func getCode(baseCode string) string {
	baseCodeParts := strings.Split(baseCode, "//")
	code := baseCodeParts[0]
	code = strings.Replace(code, " ", "", -1)
	return code
}

func getCodeType(code string) CodeType {
	if len(code) == 0 {
		return CodeTypeCommentOrEmpty
	}
	if code[0] == '(' && code[len(code)-1] == ')' {
		return CodeTypeJumpLabel
	}
	if code[0] == '@' {
		return CodeTypeInsA
	}

	if strings.Count(code, "=") == 1 || strings.Count(code, ";") == 1 {
		return CodeTypeInsC
	}
	return CodeTypeUnknow
}

// Symbols table
var Symbols map[string]string

func init() {
	Symbols = make(map[string]string)
}

// AddSymbol add a symbol to symbol table
func AddSymbol(name, value string) {
	Symbols[name] = value
}

// GetSymbol get symbol from symbol table
func GetSymbol(name string) (string, bool) {
	val, has := PreDefineSymbols[name]
	if has {
		return val, true
	}
	val, has = Symbols[name]
	return val, has
}

var addressCnt = 16

// AddVariable add a variable to symbol table
func AddVariable(name string) string {
	varValue := strconv.Itoa(addressCnt)
	AddSymbol(name, varValue)
	addressCnt++
	return varValue
}

// ParseFile scan file at first time
// at this step this function will find all jump lables
// and parse file to codes
func ParseFile(file []byte) ([]*Code, error) {
	lineNum := 0
	src := string(file)
	codeLines := strings.Split(src, "\r\n")
	parsedCodes := make([]*Code, 0, len(codeLines))
	for sourceLineNum := 0; sourceLineNum < len(codeLines); sourceLineNum++ {
		code := codeLines[sourceLineNum]
		code = getCode(code)
		codeType := getCodeType(code)
		if codeType == CodeTypeUnknow {
			return nil, fmt.Errorf("syntax error on line %d", sourceLineNum+1)
		}
		parsedCode := &Code{
			SourceLineNum: sourceLineNum,
			LineNum:       lineNum,
			Source:        code,
			Type:          codeType,
		}
		parsedCodes = append(parsedCodes, parsedCode)
		if codeType == CodeTypeInsA || codeType == CodeTypeInsC {
			lineNum++
		}
		if codeType == CodeTypeJumpLabel {
			AddSymbol(getJumpLable(code), strconv.Itoa(lineNum))
		}
	}
	return parsedCodes, nil
}

// ReplaceSymbols replace all symbols
func ReplaceSymbols(parsedCodes []*Code) string {
	lessAsmSource := ""
	instructionCnt := 0
	for i := 0; i < len(parsedCodes); i++ {
		code := parsedCodes[i]
		switch code.Type {
		case CodeTypeInsA:
			varName := code.Source[1:]
			if _, err := strconv.Atoi(varName); err != nil {
				varVal, has := GetSymbol(varName)
				if !has {
					varVal = AddVariable(varName)
				}
				code.Source = "@" + varVal
			}
			if instructionCnt != 0 {
				lessAsmSource += "\n"
			}
			lessAsmSource += code.Source
			instructionCnt++
		case CodeTypeInsC:
			if instructionCnt != 0 {
				lessAsmSource += "\n"
			}
			lessAsmSource += code.Source
			instructionCnt++
		}
	}
	return lessAsmSource
}

func formatAInstruction(old string) string {
	s := "000000000000000"
	return "0" + s[:15-len(old)] + old
}

// Parser2Binary parse to binary code
func Parser2Binary(noSymbolCodes []*Code) (string, error) {
	binaryCode := ""
	for i := 0; i < len(noSymbolCodes); i++ {
		code := noSymbolCodes[i]
		switch code.Type {
		case CodeTypeInsA:
			addressStr := code.Source[1:]
			address, err := strconv.Atoi(addressStr)
			if err != nil {
				return "", fmt.Errorf("address error in line: %d\n  @address address is invalid", code.SourceLineNum)
			}
			if float64(address) > math.Pow(2, 15)-1 || address < 0 {
				return "", fmt.Errorf("address error in line: %d\n  @address address should limit in 0~2^15 - 1", code.SourceLineNum)
			}
			code.Source = formatAInstruction(strconv.FormatInt(int64(address), 2))
			binaryCode += code.Source + "\n"
		case CodeTypeInsC:
			eqIndex := strings.Index(code.Source, "=")
			semIndex := strings.Index(code.Source, ";")
			dest := "NULL"
			comp := ""
			jump := "NULL"

			if eqIndex != -1 {
				dest = code.Source[:eqIndex]
			}
			if semIndex != -1 {
				jump = code.Source[semIndex+1:]
			}
			if eqIndex != -1 && semIndex != -1 {
				comp = code.Source[eqIndex+1 : semIndex]
			}
			if eqIndex != -1 && semIndex == -1 {
				comp = code.Source[eqIndex+1:]
			}
			if eqIndex == -1 && semIndex != -1 {
				comp = code.Source[:semIndex]
			}
			code.Source = "111" + CompCodeTable[comp] + DestCodeTable[dest] + JumpCodeTable[jump]
			binaryCode += code.Source + "\n"
			// fmt.Println(comp, ":", CompCodeTable[comp])
			// fmt.Println("111" + CompCodeTable[comp] + DestCodeTable[dest] + JumpCodeTable[jump])
		}
	}
	if strings.HasSuffix(binaryCode, "\n") {
		binaryCode = strings.TrimSuffix(binaryCode, "\n")
	}
	return binaryCode, nil
}
