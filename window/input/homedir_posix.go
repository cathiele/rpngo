//go:build !pico && !pico2

package input

import "os"

func homeDir() (string, error) {
	return os.UserHomeDir()
}
