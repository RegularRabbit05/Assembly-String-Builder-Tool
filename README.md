# Assembly String builder tool
A program to create assembly 8086 strings to print without using any printing/strings related function but only mov-xchg-int and loops
<br>
This software is probably useless in any use case so you should ignore it as it is only a fun project for school and for making assembly project easily using only the instructions we studied yet

## Sample program generated by the software that writes "Hello World!"

```x86asm
    mov DS:[ 0 ],0    ;EMPTY
    mov DS:[ 1 ],0    ;EMPTY     
    mov DS:[ 2 ], 11     ;NUMBER OF CHARS
    mov DS:[ 3 ],0    ;CH = 0
    mov DS:[ 4 ], 8     ;START OF CHARS
    mov DS:[ 5 ],0    ;BH = 0
    mov DS:[ 6 ],0    ;DL = 0
    mov DS:[ 7 ],0    ;DH = 0

    mov DS:[ 0 ],AX   ;SAVE AX AND EMPTY IT
    xchg DS:[ 2 ],CX  ;SAVE CX AND LOAD CHAR COUNTER
    xchg DS:[ 4 ],BX  ;SAVE AND LOAD POINTER
    xchg DS:[ 6 ],DX  ;SAVE AND LOAD CHARACTER

    mov DS:[ 8 ], 104  	 ;MOVE CHARS
    mov DS:[ 9 ], 101  	 ;MOVE CHARS
    mov DS:[ 10 ], 108  	 ;MOVE CHARS
    mov DS:[ 11 ], 108  	 ;MOVE CHARS
    mov DS:[ 12 ], 111  	 ;MOVE CHARS
    mov DS:[ 13 ], 32  	 ;MOVE CHARS
    mov DS:[ 14 ], 119  	 ;MOVE CHARS
    mov DS:[ 15 ], 111  	 ;MOVE CHARS
    mov DS:[ 16 ], 114  	 ;MOVE CHARS
    mov DS:[ 17 ], 108  	 ;MOVE CHARS
    mov DS:[ 18 ], 100  	 ;MOVE CHARS

    mov AX,0        ;EMPTY AX
    mov AH,2        ;OUTPUT MODE
loopPrint:
    mov DL,DS:[BX]  ;LOAD CURRENT CHAR IN MEMORY
    int 21h         ;PRINT
    inc BX          ;INCREMENT POINTER
    loop loopPrint

    mov AX,DS:[ 0 ]
    mov CX,DS:[ 2 ]
    mov BX,DS:[ 4 ]
    mov DX,DS:[ 6 ]
```
