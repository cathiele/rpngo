//go:build pico || pico2

package fileops

import "mattwach/rpngo/rpn"

const ShellHelp = ""

func Shell(r *rpn.RPN) error {
	return rpn.ErrNotSupported
}
