package main

import (
	"bufio"
	"os"
	"strconv"
)

var vmSourcrFile *os.File

// OpenVMSourcrFile open vm file return a scaner
func OpenVMSourcrFile(fileName string) (scaner *bufio.Scanner, err error) {

	vmSourcrFile, err = os.Open(fileName)
	if err != nil {
		vmSourcrFile.Close()
		return nil, err
	}
	scaner = bufio.NewScanner(vmSourcrFile)
	return
}

// CloseVMSourcrFile close vm file
func CloseVMSourcrFile() {
	if vmSourcrFile != nil {
		vmSourcrFile.Close()
	}
}

// CommandType vm command type
type CommandType uint8

const (
	CArithmetic CommandType = iota + 1
	CPush
	CPop
	CLabel
	CGoto
	CIf
	CFunction
	CReturn
	CCall
	CEmpty
)

// Command vm command
type Command struct {
	Name      string
	FirstArg  string
	SecondArg int
	Type      CommandType
}

func (command *Command) parse(commandStr string) *Command {
	commandStr += " "
	fields := []string{}
	a, b := 0, 0
	look := false
	for i := 0; i < len(commandStr); i++ {
		if !look && commandStr[i] != ' ' && commandStr[i] != '\t' {
			a = i
			look = true
			continue
		}
		if look && (commandStr[i] == ' ' || commandStr[i] == '/' || i == len(commandStr)-1) {
			b = i
			look = false
			if commandStr[a:b] == "/" {
				break
			}
			fields = append(fields, commandStr[a:b])
		}
	}
	if len(fields) == 0 {
		command.Type = CEmpty
		command.Name = "-"
		return command
	}
	command.Name = fields[0]
	command.Type = getCommandType(command.Name)
	switch command.Type {
	case CPush, CPop, CFunction, CCall:
		command.FirstArg = fields[1]
		command.SecondArg, _ = strconv.Atoi(fields[2])
	case CLabel, CIf, CGoto:
		command.FirstArg = fields[1]
	case CArithmetic:
		command.FirstArg = command.Name
	case CReturn:
	}
	return command
}

// CurrentCommand .
var CurrentCommand *Command

// HasMoreCommands check if has more commands
func HasMoreCommands(fileScaner *bufio.Scanner) bool {
	return fileScaner.Scan()
}

// Advance set new command to current command
func Advance(fileScaner *bufio.Scanner) {
	CurrentCommand = new(Command)
	CurrentCommand.parse(fileScaner.Text())
}

// GetCommandType .
func getCommandType(commandName string) CommandType {
	switch commandName {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return CArithmetic
	case "label":
		return CLabel
	case "goto":
		return CGoto
	case "if-goto":
		return CIf
	case "function":
		return CFunction
	case "call":
		return CCall
	case "push":
		return CPush
	case "pop":
		return CPop
	case "return":
		return CReturn
	}
	return 0
}
