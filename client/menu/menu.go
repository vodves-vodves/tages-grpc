package menu

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func PrintMenu(a []string) {
	fmt.Println("Menu: ")
	for i := 1; i < len(a)+1; i++ {
		if i < 10 {
			fmt.Printf("  %d » %s\n", i, a[i-1])
		}
	}
}

func UserInputInteger(format string) int {
	reader := bufio.NewReader(os.Stdin)
	var out string
	fmt.Print(format + " » ")
	out, err := reader.ReadString('\n')

	if err != nil {
		log.Println("error reading: ", err)
		return 4 // exit
	}
	if out == "" || out == "\n" {
		return 0
	}
	out = strings.TrimSuffix(out, "\r\n")
	out = strings.TrimSuffix(out, "\n")
	i, err := strconv.Atoi(out)
	if err != nil {
		log.Println("error convert: ", err)
		return 4 // exit
	}
	return i
}

func UserInput(format string) string {
	reader := bufio.NewReader(os.Stdin)
	var out string
	fmt.Print(format + " » ")
	out, err := reader.ReadString('\n')
	if err != nil {
		log.Println("error reading: ", err)
		return "" // exit
	}
	out = strings.TrimSuffix(out, "\r\n")
	out = strings.TrimSuffix(out, "\n")
	return out
}
