package io

func print(txtd TextDisplay, msg string) {
	txtd.Write([]byte(msg))
}
