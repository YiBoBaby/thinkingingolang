package mongodb

import "fmt"

func printGreen(str string) {
	fmt.Printf("%c[%dm%s%c[0m \n", 0x1b, 32, str, 0x1b)
}

func printBlue(str string) {
	fmt.Printf("%c[%dm%s%c[0m \n", 0x1b, 36, str, 0x1b)
}

func printPink(str string) {
	fmt.Printf("%c[%dm%s%c[0m \n", 0x1b, 35, str, 0x1b)
}
