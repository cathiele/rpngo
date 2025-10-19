package plotwin

import (
	"mattwach/rpngo/rpn"
)

const maxSteps = 50000

func (pw *plotWindowCommon) setProp(name string, val rpn.Frame) error {
	switch name {
	case "autosteps":
		v, err := val.BoundedInt(1, maxSteps)
		if err != nil {
			return err
		}
		pw.autosteps = uint32(v)
	case "minx":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.minx = v
		if pw.maxx <= pw.minx {
			pw.maxx = pw.minx + 1
		}
		pw.autox = false
	case "maxx":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.maxx = v
		if pw.maxx <= pw.minx {
			pw.minx = pw.maxx - 1
		}
		pw.autox = false
	case "miny":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.miny = v
		if pw.maxy <= pw.miny {
			pw.maxy = pw.miny + 1
		}
		pw.autoy = false
	case "maxy":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.maxy = v
		if pw.maxy <= pw.miny {
			pw.miny = pw.maxy - 1
		}
		pw.autoy = false
	case "minv":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.minv = v
		if pw.maxv <= pw.minv {
			pw.maxv = pw.minv + 1
		}
	case "maxv":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.maxv = v
		if pw.maxv <= pw.minv {
			pw.minv = pw.maxv - 1
		}
	case "autox":
		v, err := val.Bool()
		if err != nil {
			return err
		}
		pw.autox = v
	case "autoy":
		v, err := val.Bool()
		if err != nil {
			return err
		}
		pw.autoy = v
	case "steps":
		v, err := val.BoundedInt(1, maxSteps)
		if err != nil {
			return err
		}
		pw.steps = uint32(v)
	default:
		return rpn.ErrUnknownProperty
	}

	return nil
}

func (pw *plotWindowCommon) getProp(name string) (rpn.Frame, error) {
	switch name {
	case "autosteps":
		return rpn.IntFrame(int64(pw.autosteps), rpn.INTEGER_FRAME), nil
	case "minx":
		return rpn.RealFrame(pw.minx), nil
	case "maxx":
		return rpn.RealFrame(pw.maxx), nil
	case "miny":
		return rpn.RealFrame(pw.miny), nil
	case "maxy":
		return rpn.RealFrame(pw.maxy), nil
	case "minv":
		return rpn.RealFrame(pw.minv), nil
	case "maxv":
		return rpn.RealFrame(pw.maxv), nil
	case "autox":
		return rpn.BoolFrame(pw.autox), nil
	case "autoy":
		return rpn.BoolFrame(pw.autoy), nil
	case "steps":
		return rpn.IntFrame(int64(pw.steps), rpn.INTEGER_FRAME), nil
	}
	return rpn.Frame{}, rpn.ErrUnknownProperty
}

var props = []string{"autox", "autoy", "minv", "maxv", "minx", "maxx", "miny", "maxy", "steps"}

func (pw *plotWindowCommon) ListProps() []string {
	return props
}
