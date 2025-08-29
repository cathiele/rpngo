package rpn

import "strings"

func (r *RPN) PushMessage(msg string) {
	r.messages = append(r.messages, msg)
}

func (r *RPN) PopMessages() string {
	if len(r.messages) == 0 {
		return ""
	}
	msgs := strings.Join(r.messages, "\n")
	r.messages = r.messages[:0]
	return msgs
}
