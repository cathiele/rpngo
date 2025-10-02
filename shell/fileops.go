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

const SaveHelp = "Uses $0 as the filename to save the contents of $1  Both are popped from the stack.  A '\\n' character is added.  Use append if this is not wanted."

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
	return os.WriteFile(path, []byte(data.String(false)+"\n"), 0644)
}

const AppendHelp = "Uses $0 as the filename to append the contents of $1  Both are popped from the stack.  No '\\n' character is added."

func Append(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	data, err := r.PopFrame()
	if err != nil {
		r.PushString(path)
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(data.String(false)))
	return err
}

const ChangeDirHelp = "Change the working directory"

func ChangeDir(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	return os.Chdir(path)
}
