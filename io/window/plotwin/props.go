package plotwin

import (
	"errors"
	"fmt"
	"mattwach/rpngo/rpn"
)

const maxSteps = 50000

func (pw *PlotWindow) SetProp(name string, val rpn.Frame) error {
	switch name {
	case "minx":
		if val.Type != rpn.COMPLEX_FRAME {
			return rpn.ErrExpectedANumber
		}
		pw.minx = real(val.Complex)
		if pw.maxx <= pw.minx {
			pw.maxx = pw.minx + 1
		}
		pw.autox = false
	case "maxx":
		if val.Type != rpn.COMPLEX_FRAME {
			return rpn.ErrExpectedANumber
		}
		pw.maxx = real(val.Complex)
		if pw.maxx <= pw.minx {
			pw.minx = pw.maxx - 1
		}
		pw.autox = false
	case "miny":
		if val.Type != rpn.COMPLEX_FRAME {
			return rpn.ErrExpectedANumber
		}
		pw.miny = real(val.Complex)
		if pw.maxy <= pw.miny {
			pw.maxy = pw.miny + 1
		}
		pw.autoy = false
	case "maxy":
		if val.Type != rpn.COMPLEX_FRAME {
			return rpn.ErrExpectedANumber
		}
		pw.maxy = real(val.Complex)
		if pw.maxy <= pw.miny {
			pw.miny = pw.maxy - 1
		}
		pw.autoy = false
	case "minv":
		if val.Type != rpn.COMPLEX_FRAME {
			return rpn.ErrExpectedANumber
		}
		pw.minv = real(val.Complex)
		if pw.maxv <= pw.minv {
			pw.maxv = pw.minv + 1
		}
	case "maxv":
		if val.Type != rpn.COMPLEX_FRAME {
			return rpn.ErrExpectedANumber
		}
		pw.maxv = real(val.Complex)
		if pw.maxv <= pw.minv {
			pw.minv = pw.maxv - 1
		}
	case "autox":
		if val.Type != rpn.BOOL_FRAME {
			return rpn.ErrExpectedABoolean
		}
		pw.autox = val.Bool
	case "autoy":
		if val.Type != rpn.BOOL_FRAME {
			return rpn.ErrExpectedABoolean
		}
		pw.autoy = val.Bool
	case "steps":
		if val.Type != rpn.COMPLEX_FRAME {
			return rpn.ErrExpectedANumber
		}
		v := real(val.Complex)
		if v < 1 {
			return errors.New("steps must be >= 1")
		}
		if v > maxSteps {
			return fmt.Errorf("steps must be <= %d", maxSteps)
		}
		pw.steps = uint32(v)
	default:
		return fmt.Errorf("unknown property: %s", name)
	}

	return nil
}

func (pw *PlotWindow) GetProp(name string) (rpn.Frame, error) {
	switch name {
	case "minx":
		return rpn.Frame{Type: rpn.COMPLEX_FRAME, Complex: complex(pw.minx, 0)}, nil
	case "maxx":
		return rpn.Frame{Type: rpn.COMPLEX_FRAME, Complex: complex(pw.maxx, 0)}, nil
	case "miny":
		return rpn.Frame{Type: rpn.COMPLEX_FRAME, Complex: complex(pw.miny, 0)}, nil
	case "maxy":
		return rpn.Frame{Type: rpn.COMPLEX_FRAME, Complex: complex(pw.maxy, 0)}, nil
	case "minv":
		return rpn.Frame{Type: rpn.COMPLEX_FRAME, Complex: complex(pw.minv, 0)}, nil
	case "maxv":
		return rpn.Frame{Type: rpn.COMPLEX_FRAME, Complex: complex(pw.maxv, 0)}, nil
	case "autox":
		return rpn.Frame{Type: rpn.BOOL_FRAME, Bool: pw.autox}, nil
	case "autoy":
		return rpn.Frame{Type: rpn.BOOL_FRAME, Bool: pw.autoy}, nil
	case "steps":
		return rpn.Frame{Type: rpn.COMPLEX_FRAME, Complex: complex(float64(pw.steps), 0)}, nil
	}
	return rpn.Frame{}, fmt.Errorf("unknown property: %s", name)
}

func (pw *PlotWindow) ListProps() []string {
	return []string{"autox", "autoy", "minv", "maxv", "minx", "maxx", "miny", "maxy", "steps"}
}
