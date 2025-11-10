package input

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"strings"
)

const MAX_SHOW_FRAMES = 1000

func (iw *InputWindow) SetProp(name string, val rpn.Frame) error {
	switch name {
	case "autofn":
		var fn []string
		addfn := func(s string) error {
			fn = append(fn, s)
			return nil
		}
		if err := parse.Fields(val.String(false), addfn); err != nil {
			return err
		}
		iw.autofn = fn
		return nil
	case "histpath":
		if !val.IsString() {
			return rpn.ErrExpectedAString
		}
		iw.gl.histpath = val.UnsafeString()
		return nil
	case "showframes":
		v, err := val.BoundedInt(0, MAX_SHOW_FRAMES)
		if err != nil {
			return err
		}
		iw.showFrames = int(v)
		return nil
	}
	return rpn.ErrNotSupported
}

func (iw *InputWindow) GetProp(name string) (rpn.Frame, error) {
	switch name {
	case "autofn":
		return rpn.StringFrame(strings.Join(iw.autofn, " "), rpn.STRING_BRACES), nil
	case "histpath":
		return rpn.StringFrame(iw.gl.histpath, rpn.STRING_SINGLE_QUOTE), nil
	case "showframes":
		return rpn.IntFrame(int64(iw.showFrames), rpn.INTEGER_FRAME), nil
	}
	return rpn.Frame{}, rpn.ErrNotSupported
}

var inputProps = []string{"autofn", "histpath", "showframes"}

func (iw *InputWindow) ListProps() []string {
	return inputProps
}
