package key

// ASCII for most keys, except the special ones below
type Key int

const (
	KEY_LEFT Key = iota + 256
	KEY_RIGHT
	KEY_BACKSPACE
)
