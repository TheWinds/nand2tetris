// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.
@color
M=0

(MAINLOOP)
@color
M=0
@SCREEN
D=A
@i
M=D

@KBD
D=M
@FILLSCREEN
D;JEQ
@SETBLACK
0;JMP

(SETBLACK)
@color
M=-1
@FILLSCREEN
0;JMP

(FILLSCREEN)
@i
D=M
@KBD
D=D-A
@MAINLOOP
D;JEQ
@color
D=M
@i
A=M
M=D
@i
M=M+1
@FILLSCREEN
0;JMP