// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"fmt"
	"machine"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"time"
)

func main() {
	time.Sleep(2 * time.Second)
	var r rpn.RPN  // object allocated on the heap: escapes at line 29
	r.Init()
	functions.RegisterAll(&r)

	fmt.Println("Type ? for help or topic? for more detailed help")

	for {
		line := readLine()
		args, err := parse.Fields(line)

		if err == nil {
			err = r.Exec(args)
		}

		if err == nil {
			r.IterFrames(func(sf rpn.Frame) {
				fmt.Println(sf.String(true))  // object allocated on the heap: escapes at line 34
			})
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func readLine() string {
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
				return string(msg)
			default:
				msg = append(msg, c)
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}