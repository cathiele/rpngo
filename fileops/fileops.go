package fileops

import (
	"fmt"
	"mattwach/rpngo/rpn"
)

type FileOpsDriver interface {
	FileSize(path string) (int, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	AppendToFile(path string, data []byte) error
	Chdir(path string) error
}

type FileOps struct {
	maxFileSize int
	driver      FileOpsDriver
}

func (fo *FileOps) InitAndRegister(r *rpn.RPN, maxFileSize int, shellAvailable bool, driver FileOpsDriver) {
	fo.maxFileSize = maxFileSize
	fo.driver = driver
	fo.register(r, shellAvailable)
}

const LoadHelp = "Loads the given filename and places it on the stack as a string variable"

func (fo *FileOps) Load(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	sz, err := fo.driver.FileSize(path)
	if err != nil {
		return err
	}
	if sz > fo.maxFileSize {
		return fmt.Errorf("file is too large.  %v > %v max bytes", sz, fo.maxFileSize)
	}
	data, err := fo.driver.ReadFile(path)
	if err != nil {
		return err
	}
	return r.PushString(string(data))
}

const SaveHelp = "Uses $0 as the filename to save the contents of $1  Both are popped from the stack.  A '\\n' character is added.  Use append if this is not wanted."

func (fo *FileOps) Save(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	data, err := r.PopFrame()
	if err != nil {
		r.PushString(path)
		return err
	}
	return fo.driver.WriteFile(path, []byte(data.String(false)+"\n"))
}

const AppendHelp = "Uses $0 as the filename to append the contents of $1  Both are popped from the stack.  No '\\n' character is added."

func (fo *FileOps) Append(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	data, err := r.PopFrame()
	if err != nil {
		r.PushString(path)
		return err
	}
	return fo.driver.AppendToFile(path, []byte(data.String(false)))
}

const ChangeDirHelp = "Change the working directory"

func (fo *FileOps) ChangeDir(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	return fo.driver.Chdir(path)
}
