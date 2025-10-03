// Package key defines key presses.
package key

// ASCII for most keys, except the special ones below
type Key int

const (
	KEY_LEFT Key = iota + 256
	KEY_RIGHT
	KEY_UP
	KEY_DOWN
	KEY_BACKSPACE
	KEY_DEL
	KEY_INS
	KEY_END
	KEY_HOME
	KEY_EOF
	KEY_F1
	KEY_F2
	KEY_F3
	KEY_F4
	KEY_F5
	KEY_F6
	KEY_F7
	KEY_F8
	KEY_F9
	KEY_F10
	KEY_F11
	KEY_F12
)
