package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestX(t *testing.T) {
	fmt.Println(strings.Split("//", "//"))
	fmt.Println("xx", strconv.FormatInt(8, 2))
}

func getTestSource() (src, lessSrc, binarySrc string) {
	src = `// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/06/max/Max.asm

// Computes R2 = max(R0, R1)  (R0,R1,R2 refer to RAM[0],RAM[1],RAM[2])

   @R0
   D=M              // D = first number
   @R1
   D=D-M            // D = first number - second number
   @OUTPUT_FIRST
   D;JGT            // if D>0 (first is greater) goto output_first
   @R1
   D=M              // D = second number
   @OUTPUT_D
   0;JMP            // goto output_d
(OUTPUT_FIRST)
   @R0             
   D=M              // D = first number
(OUTPUT_D)
   @R2
   M=D              // M[2] = D (greatest number)
(INFINITE_LOOP)
   @INFINITE_LOOP
   0;JMP            // infinite loop`

	lessSrc = `@0
D=M
@1
D=D-M
@10
D;JGT
@1
D=M
@12
0;JMP
@0
D=M
@2
M=D
@14
0;JMP`
	binarySrc = `0000000000000000
1111110000010000
0000000000000001
1111010011010000
0000000000001010
1110001100000001
0000000000000001
1111110000010000
0000000000001100
1110101010000111
0000000000000000
1111110000010000
0000000000000010
1110001100001000
0000000000001110
1110101010000111`
	return
}

func TestParser(t *testing.T) {
	src, lessSrc, binarySrc := getTestSource()
	parsedCodes, err := ParseFile([]byte(src))
	Convey("test parse file", t, func() {
		So(err, ShouldBeNil)
		So(Symbols["OUTPUT_FIRST"], ShouldEqual, "10")
		So(Symbols["OUTPUT_D"], ShouldEqual, "12")
		So(Symbols["INFINITE_LOOP"], ShouldEqual, "14")
	})
	srcLessSymbols := ReplaceSymbols(parsedCodes)
	Convey("test replace symbols", t, func() {
		So(lessSrc, ShouldEqual, srcLessSymbols)
	})
	binaryCode, err := Parser2Binary(parsedCodes)
	Convey("test parse to binary code", t, func() {
		So(err, ShouldBeNil)
		So(binarySrc, ShouldEqual, binaryCode)
	})

}
