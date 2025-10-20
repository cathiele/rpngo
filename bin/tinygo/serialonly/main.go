// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"fmt"
	"machine"
	"mattwach/rpngo/bin/tinygo"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"time"
)

var fields []string = make([]string, 128)

func main() {
	time.Sleep(2 * time.Second)
	var r rpn.RPN // object allocated on the heap (OK)
	r.Init(256)
	functions.RegisterAll(&r)

	print("Type ? for help or topic? for more detailed help")

	for {
		tinygo.DumpMemStats()
		line := readLine()
		fields = fields[:0]
		var err error
		fields, err = parse.Fields(line, fields)

		if err == nil {
			err = r.Exec(fields)
		}

		if err == nil {
			for _, f := range r.Frames {
				for _, r := range f.String(true) {
					machine.Serial.WriteByte(byte(r))
				}
				machine.Serial.WriteByte('\n')
			}
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func readLine() string {
	var msg []byte
	print("> ")
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
