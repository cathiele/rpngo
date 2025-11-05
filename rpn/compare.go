package rpn

const tolerance = 1e-9

func checkFloatEqual(a, b complex128) bool {
	d := real(a) - real(b)
	if d < -tolerance {
		return false
	}
	if d > tolerance {
		return false
	}
	d = imag(a) - imag(b)
	if d < -tolerance {
		return false
	}
	return d <= tolerance
}

// Ordering of non-comparable types
// bool < number < string

func (a Frame) IsLessThan(b Frame) bool {
	switch a.ftype {
		case COMPLEX_FRAME:
			switch b.ftype {
				case COMPLEX_FRAME:
					return real(a.cmplx) < real(b.cmplx)
				case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
					return real(a.cmplx) < float64(b.intv)
				case STRING_FRAME:
					return true
				case BOOL_FRAME:
					return false
			}
		case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
			switch b.ftype {
				case COMPLEX_FRAME:
					return float64(a.intv) < real(b.cmplx)
				case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
					return a.intv < b.intv
				case STRING_FRAME:
					return true
				case BOOL_FRAME:
					return false
			}
		case STRING_FRAME:
			switch b.ftype {
				case STRING_FRAME:
					return a.str < b.str
				default:
					return false
			}
		case BOOL_FRAME:
			switch b.ftype {
				case BOOL_FRAME:
					return a.intv < b.intv
				default:
					return true
			}
	}

	return false
}

func (a Frame) IsLessThanOrEqual(b Frame) bool {
	switch a.ftype {
		case COMPLEX_FRAME:
			switch b.ftype {
				case COMPLEX_FRAME:
					return checkFloatEqual(a.cmplx, b.cmplx) || (real(a.cmplx) < real(b.cmplx))
				case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
					return checkFloatEqual(a.cmplx, complex(float64(b.intv), 0)) || (real(a.cmplx) < float64(b.intv))
				case STRING_FRAME:
					return true
				case BOOL_FRAME:
					return false
			}
		case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
			switch b.ftype {
				case COMPLEX_FRAME:
					return checkFloatEqual(complex(float64(a.intv), 0), b.cmplx) || float64(a.intv) < real(b.cmplx)
				case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
					return a.intv <= b.intv
				case STRING_FRAME:
					return true
				case BOOL_FRAME:
					return false
			}
		case STRING_FRAME:
			switch b.ftype {
				case STRING_FRAME:
					return a.str <= b.str
				default:
					return false
			}
		case BOOL_FRAME:
			switch b.ftype {
				case BOOL_FRAME:
					return a.intv <= b.intv
				default:
					return true
			}
	}

	return false
}

func (a Frame) IsEqual(b Frame) bool {
	switch a.ftype {
		case COMPLEX_FRAME:
			switch b.ftype {
				case COMPLEX_FRAME:
					return checkFloatEqual(a.cmplx, b.cmplx)
				case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
					return checkFloatEqual(a.cmplx, complex(float64(b.intv), 0))
				case STRING_FRAME:
					return false
				case BOOL_FRAME:
					return false
			}
		case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
			switch b.ftype {
				case COMPLEX_FRAME:
					return checkFloatEqual(complex(float64(a.intv), 0), b.cmplx)
				case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
					return a.intv == b.intv
				case STRING_FRAME:
					return false
				case BOOL_FRAME:
					return false
			}
		case STRING_FRAME:
			switch b.ftype {
				case STRING_FRAME:
					return a.str == b.str
				default:
					return false
			}
		case BOOL_FRAME:
			switch b.ftype {
				case BOOL_FRAME:
					return a.intv == b.intv
				default:
					return false
			}
	}

	return false
}

