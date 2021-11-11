package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var text string
	var reader *bufio.Reader
	var code string

	reader = bufio.NewReader(os.Stdin)

	fmt.Println(cleanup())

	fmt.Println("This program lets you build an assembly print from a string (using only mov & int)")
	fmt.Println("This program will automatically generate code for copying all the used Registers into DATA Segment in order to resume execution at a later point")
	fmt.Println("but be careful as the code will NOT clean the used DATA Segment bytes after use leaving them dirty with old information, about the registers, and the string intact")
	fmt.Println("(useful in case the string is printed more than once as the code may be reused but remember to copy the 'Register backup section' of it)")
	fmt.Println("The code will also generate indentation and comments by itself about what the assembly is doing. The loop function is called 'loopPrint' so please leave this name unused!")
	fmt.Println()
	fmt.Println("-- Exit anytime using ^D/^C --")

	for {
		fmt.Println("Please input the segment mode [DATA/CODE]")
		printCursor()
		text, _ = reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		text = strings.ToLower(text)

		if strings.Compare("data", text) == 0 {
			fmt.Println("Assuming DATA segment for string storage (in case it needs to be reused)")
			code = getString(reader, false)
			break
		} else if strings.Compare("code", text) == 0 {
			fmt.Println("Assuming CODE segment for string storage (in case of single usage)")
			code = getString(reader, true)
			break
		} else {
			fmt.Println("Unknown choice...")
		}
	}

	fmt.Println("Done converting please copy the code or continue to pick if saving to file!\nPress 'Enter' to continue...")
	reader.ReadBytes('\n')

	if askTrueFalse("Would you like to save to file? [YES/NO]", reader) {
		writeToFile(code, reader)
	}

	os.Exit(0)
}

func askTrueFalse(question string, reader *bufio.Reader) bool {
	var text string
	for {
		fmt.Println(question)
		printCursor()
		text, _ = reader.ReadString('\n')
		text = string(text[0])
		text = strings.ToLower(text)

		if strings.Compare("y", text) == 0 {
			return true
		} else if strings.Compare("n", text) == 0 {
			return false
		} else {
			fmt.Println("Unknown choice...")
		}
	}
}

func printCursor() {
	fmt.Print("-> ")
}

func printConverted(data []byte) {
	fmt.Println("Converted: ", data)
}

func getString(reader *bufio.Reader, isCode bool) string {
	var text string
	var usingCarriage bool
	var toBytes []byte

	if askTrueFalse("Would you like to automatically add \"carriage return\" to \"newlines\"? [YES/NO]:", reader) {
		usingCarriage = true
		fmt.Println("Using auto carriage return")
	} else {
		usingCarriage = false
		fmt.Println("Not using auto carriage return")
	}

	fmt.Println("Please input your string to compute:")
	printCursor()
	text, _ = reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)

	if usingCarriage {
		text = strings.ReplaceAll(text, "\\r", "")
		text = strings.ReplaceAll(text, "\\n", "\n\r")
	} else {
		text = strings.ReplaceAll(text, "\\n", "\n")
		text = strings.ReplaceAll(text, "\\r", "\r")
	}

	toBytes = strToAscii(text, reader)
	printConverted(toBytes)

	if isCode {
		return codeSegCompute(toBytes, reader)
	} else {
		return dataSegCompute(toBytes, reader)
	}
}

func strToAscii(text string, reader *bufio.Reader) []byte {
	var result []byte
	var chars []rune

	chars = []rune(text)

	for i := 0; i < len(chars); i++ {
		result = append(result, byte(chars[i]))
	}

	if len(result) <= 0 {
		fmt.Println("No data was inputted, leaving...")
		fmt.Println("Press 'Enter' to continue...")
		reader.ReadBytes('\n')
		os.Exit(1)
	}

	return result
}

func cleanup() string {
	return fmt.Sprint("\n\n\n\n")
}

func dataSegCompute(data []byte, reader *bufio.Reader) string {
	var code string
	var location uint16
	var currentLoc uint16
	var length uint8
	var i int

	code = ""
	length = uint8(len(data))
	fmt.Println("Attention: this operation will take up to ", length+8, " bytes in the data segment memory")
	fmt.Println("Please input a starting memory location (DS:[YourLocation])")
	printCursor()
	if _, err := fmt.Scan(&location); err != nil {
		log.Print("Error: Unable to read location: ", err)
		fmt.Println("Press 'Enter' to continue...")
		reader.ReadBytes('\n')
		os.Exit(1)
	}

	if uint32(location)+uint32(length)+8 > 65535 {
		fmt.Println("Error: there is not enough space in the data segment to store the bytes necessary for this operation! (8086: offset+len > 0xFFFF)")
		fmt.Println("Press 'Enter' to continue...")
		reader.ReadBytes('\n')
		os.Exit(1)
	}
	fmt.Println(cleanup())

	currentLoc = location
	code += fmt.Sprintln("    mov DS:[", currentLoc, "],0    ;EMPTY\n    mov DS:[", currentLoc+1, "],0    ;EMPTY     \n    mov DS:[", currentLoc+2, "],", length, "    ;NUMBER OF CHARS\n    mov DS:[", currentLoc+3, "],0    ;CH = 0\n    mov DS:[", currentLoc+4, "],", currentLoc+8, "    ;START OF CHARS\n    mov DS:[", currentLoc+5, "],0    ;BH = 0\n    mov DS:[", currentLoc+6, "],0    ;DL = 0\n    mov DS:[", currentLoc+7, "],0    ;DH = 0")
	code += fmt.Sprintln()
	code += fmt.Sprintln("    mov DS:[", currentLoc, "],AX   ;SAVE AX AND EMPTY IT\n    xchg DS:[", currentLoc+2, "],CX  ;SAVE CX AND LOAD CHAR COUNTER\n    xchg DS:[", currentLoc+4, "],BX  ;SAVE AND LOAD POINTER\n    xchg DS:[", currentLoc+6, "],DX  ;SAVE AND LOAD CHARACTER")
	code += fmt.Sprintln()

	i = 0
	for currentLoc += 8; uint32(currentLoc) < uint32(length)+8+uint32(location); currentLoc++ {
		code += fmt.Sprintln("    mov DS:[", currentLoc, "],", data[i], " 	 ;MOVE CHARS")
		i++
	}
	code += fmt.Sprintln()
	code += fmt.Sprintln("    mov AX,0        ;EMPTY AX\n    mov AH,2        ;OUTPUT MODE")
	code += fmt.Sprintln("loopPrint:\n    mov DL,DS:[BX]  ;LOAD CURRENT CHAR IN MEMORY\n    int 21h         ;PRINT\n    inc BX          ;INCREMENT POINTER\n    loop loopPrint")
	code += fmt.Sprintln()
	code += fmt.Sprintln("    mov AX,DS:[", location, "]\n    mov CX,DS:[", location+2, "]\n    mov BX,DS:[", location+4, "]\n    mov DX,DS:[", location+6, "]")

	fmt.Println(code)
	fmt.Println(cleanup())
	return code
}

func codeSegCompute(data []byte, reader *bufio.Reader) string {
	var code string
	var location uint16
	var i int

	code = ""
	fmt.Println("Please input a memory location for a backup of AX and DX in order to resume execution at a later point, the code will need 4 bytes to work (DS:[YourLocation])")
	printCursor()
	if _, err := fmt.Scan(&location); err != nil {
		log.Print("Error: Unable to read location: ", err)
		fmt.Println("Press 'Enter' to continue...")
		reader.ReadBytes('\n')
		os.Exit(1)
	}

	if uint32(location)+4 > 65535 {
		fmt.Println("Error: there is not enough space in the data segment to store the bytes necessary for this operation! (8086: 4+len > 0xFFFF)")
		os.Exit(1)
	}

	fmt.Println(cleanup())

	code += fmt.Sprintln("    mov DS:[", location, "],AX   ;SAVE AX\n    mov DS:[", location+2, "],DX   ;SAVE DX\n    mov AH,2        ;PREPARE FOR OUTPUT")
	code += fmt.Sprintln()
	for i = 0; i < len(data); i++ {
		code += fmt.Sprintln("    mov DL, ", data[i], " ;LOAD CHARACTER\n    int 21h     ;BIOS INTERRUPT FOR PRINTING")
	}
	code += fmt.Sprintln()
	code += fmt.Sprintln("    mov DX,DS:[", location+2, "]\n    mov AX,DS:[", location, "]")

	fmt.Println(code)
	fmt.Println(cleanup())
	return code
}

func writeToFile(code string, reader *bufio.Reader) {
	var bytes []byte
	var text string

	bytes = []byte(code)

	fmt.Println("Write a name for the source file (also supply an extension, usually .s or .asm)")
	printCursor()
	text, _ = reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)

	file, err1 := os.Create(text)

	if err1 != nil {
		fmt.Println("Error: unable to create the assembly file please save the code above, the program will exit: ", err1)
		fmt.Println("Press 'Enter' to continue...")
		reader.ReadBytes('\n')
		os.Exit(1)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	_, err2 := file.Write(bytes)

	if err2 != nil {
		fmt.Println("Error: unable to write to file please save the code above, the program will exit: ", err2)
		fmt.Println("Press 'Enter' to continue...")
		reader.ReadBytes('\n')
		os.Exit(1)
	}

	fmt.Println("The file has been saved to '", text, "'!")
}
