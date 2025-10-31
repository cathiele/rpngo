// RPNStack holds stack information for an RPNCalc
package rpn

func (r *RPN) Clear() {
	r.Frames = r.Frames[:0]
}

func (r *RPN) StackLen() int {
	return len(r.Frames)
}

func (r *RPN) PushFrame(f Frame) error {
	if len(r.Frames) >= r.maxStackDepth {
		return ErrStackFull
	}
	r.Frames = append(r.Frames, f)
	return nil
}

func (r *RPN) PopFrame() (sf Frame, err error) {
	if len(r.Frames) == 0 {
		err = ErrStackEmpty
		return
	}
	sf = r.Frames[len(r.Frames)-1]
	r.Frames = r.Frames[:len(r.Frames)-1]
	return
}

func (r *RPN) Pop2Frames() (a Frame, b Frame, err error) {
	b, err = r.PopFrame()
	if err != nil {
		return
	}
	a, err = r.PopFrame()
	return
}

func (r *RPN) PeekFrame(framesBack int) (sf Frame, err error) {
	if framesBack < 0 {
		err = ErrIllegalValue
		return
	}
	if framesBack >= len(r.Frames) {
		err = ErrNotEnoughStackFrames
		return
	}
	sf = r.Frames[len(r.Frames)-1-framesBack]
	return
}

func (r *RPN) DeleteFrame(framesBack int) (sf Frame, err error) {
	sf, err = r.PeekFrame(framesBack)
	if err != nil {
		return
	}
	idx := len(r.Frames) - 1 - framesBack
	copy(r.Frames[idx:], r.Frames[idx+1:])
	r.Frames = r.Frames[:len(r.Frames)-1]
	return
}

func (r *RPN) InsertFrame(f Frame, framesBack int) error {
	if framesBack < 0 {
		return ErrIllegalValue
	}
	if framesBack > len(r.Frames) {
		return ErrNotEnoughStackFrames
	}
	if framesBack == 0 {
		return r.PushFrame(f)
	}
	idx := len(r.Frames) - framesBack
	r.Frames = append(r.Frames, Frame{})
	copy(r.Frames[idx+1:], r.Frames[idx:])
	r.Frames[idx] = f
	return nil
}

const stackSizeHelp = "Pushes the current stack size to the stack (non-inclusive)."

func stackSize(r *RPN) error {
	return r.PushFrame(IntFrame(int64(len(r.Frames)), INTEGER_FRAME))
}
