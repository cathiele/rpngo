package io

func getLine(input Input, txtd TextDisplay) (string, error) {
	print(txtd, "> ")
	var line []byte
	for {
		b, err := input.ReadByte()
		if err != nil {
			return "", err
		}
		puts(txtd, b)
		if b == '\n' {
			break
		}
		line = append(line, b)
	}
	return string(line), nil
}
