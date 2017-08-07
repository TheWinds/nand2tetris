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
// 保存颜色信息
@color
M=0
// 主循环,不断监听键盘
(MAINLOOP)
@color
M=0
@SCREEN
D=A
// addr 屏幕当前被填充的地址
@addr
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
@addr
D=M
// 判断是否绘制完成
@KBD
D=D-A
@MAINLOOP
D;JEQ
// 绘制
@color
D=M
@addr
A=M
M=D
// 到下一个屏幕填充的地址
@addr
M=M+1
@FILLSCREEN
0;JMP