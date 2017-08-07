// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)

// Put your code here.
// init variables
@R0
D=M
@base
M=D
@R1
D=M
@n
M=D
@i
M=1
@result
M=0
@R2
M=0
// mult
(LOOP)
@i
D=M
@n
D=D-M
@OVER
D;JGT
@base
D=M
@result
D=D+M
M=D
@i
M=M+1
@LOOP
0;JMP
(OVER)
@result
D=M
@R2
M=D
@END
0;JMP
(END)
@END
0;JMP