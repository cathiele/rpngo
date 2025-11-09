//go:build pico || pico2

package input

func homeDir() (string, error) {
	return "/", nil
}
