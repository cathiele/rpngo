// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"fmt"
	"machine"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/rpn"
	"strings"
	"time"
)

func main() {
	time.Sleep(2 * time.Second)
	var r rpn.RPN
	r.Init()
	functions.RegisterAll(&r)

	fmt.Println("Type ? for help or topic? for more detailed help")

	for {
		args := readLine()

		err := r.Exec(args)

		if err == nil {
			r.IterFrames(func(sf rpn.Frame) {
				fmt.Println(sf.String(true))
			})
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func readLine() []string {
	var msg []byte
	fmt.Print("> ")
	for {
		c, err := machine.Serial.ReadByte()
		machine.Serial.WriteByte(c)
		if err == nil {
			switch c {
			case 0:
				break
			case 13:
				machine.Serial.WriteByte('\n')
				return strings.Fields(string(msg))
			default:
				msg = append(msg, c)
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}
