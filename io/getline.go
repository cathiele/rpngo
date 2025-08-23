package io

func getLine(input Input, txtd TextDisplay) (string, error) {
	print(txtd, "> ")
	var line []byte
	for {
		b, err := input.ReadByte()
		if err != nil {
			return "", err
		}
		txtd.Write([]byte{b})
		line = append(line, b)
		if b == '\n' {
			break
		}
	}
	return string(line), nil
}
