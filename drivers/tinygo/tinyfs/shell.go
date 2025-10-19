//go:build pico || pico2

package tinyfs

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	errArgsNotSupported       = errors.New("args not supported")
	errEmptyArgument          = errors.New("empty argument")
	errOnlyOnePathIsSupported = errors.New("only one path is supported")
	errUnknownCommand         = errors.New("unknown command")
)

// This file implements a very basic shell that implements a minimal
// interface of what could be considered "necessary" commands.
// It's done as shell emulation so that the PC and PicoCalc
// operations will use similar syntax (althogh the PC version is
// much mor robust with a real shell, of course).

func (fo *FileOpsDriver) Shell(args []string, stdin io.Reader) (string, error) {
	if fo.initErr != nil {
		return "", fo.initErr
	}
	if len(args) == 0 {
		return "", nil
	}
	var val string
	var err error
	switch args[0] {
	case "ls":
		val, err = fo.ls(args[1:])
	case "pwd":
		val, err = fo.getpwd(args[1:])
	default:
		return "", errUnknownCommand
	}

	return val, err
}

func (fo *FileOpsDriver) ls(args []string) (string, error) {
	path := ""
	longMode := false
	for _, arg := range args {
		if arg == "" {
			return "", errEmptyArgument
		} else if arg[0] != '-' {
			if len(path) > 0 {
				return "", errOnlyOnePathIsSupported
			}
			path = arg
		} else if arg == "-l" {
			longMode = true
		} else {
			return "", fmt.Errorf("unknown flag: %v", arg)
		}
	}
	path = absPath(fo.pwd, path, true, false) // object allocated on the heap: escapes at line 65
	print("open ")
	println(path)
	dir, err := fo.fs.Open(path)
	if err != nil {
		return "", err
	}
	defer dir.Close()
	println("read dir")
	infos, err := dir.Readdir(0)
	if err != nil {
		return "", err
	}
	println("for loop")
	buff := make([]byte, 0, 128) // object allocated on the heap: escapes at line 77
	for _, info := range infos {
		if longMode {
			buff = appendLongInfo(buff, info)
		} else {
			buff = appendShortInfo(buff, info)
		}
	}
	println("done")
	return string(buff), nil
}

func appendLongInfo(buff []byte, info os.FileInfo) []byte {
	if info.IsDir() {
		buff = append(buff, byte('d'))
	} else {
		buff = append(buff, byte('-'))
	}
	return append(buff, []byte(fmt.Sprintf("rwxrwxrwx %5d %s\n", info.Size(), info.Name()))...) // object allocated on the heap: escapes at line 96
}

func appendShortInfo(buff []byte, info os.FileInfo) []byte {
	buff = append(buff, []byte(info.Name())...)
	if info.IsDir() {
		buff = append(buff, byte('/'))
	}
	return append(buff, byte('\n'))
}

func (fo *FileOpsDriver) getpwd(args []string) (string, error) {
	if len(args) != 0 {
		return "", errArgsNotSupported
	}
	return "/" + fo.pwd + "\n", nil
}
