package shell

import (
	"fmt"
	"mattwach/rpngo/rpn"
	"os"
)

const maxFileSize = 65536

const LoadHelp = "Loads the given filename and places it on the stack as a string variable"

func Load(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.Size() > maxFileSize {
		return fmt.Errorf("file is too large.  %v > %v max bytes", s.Size(), maxFileSize)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return r.PushString(string(data))
}

const SaveHelp = "Uses $0 as the filename to save the contents of $1  Both are popped from the stack."

func Save(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	data, err := r.PopFrame()
	if err != nil {
		r.PushString(path)
		return err
	}
	return os.WriteFile(path, []byte(data.String(false)+"\n"), 0666)
}

const ChangeDirHelp = "Change the working directory"

func ChangeDir(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	return os.Chdir(path)
}
