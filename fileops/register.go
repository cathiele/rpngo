package fileops

import "mattwach/rpngo/rpn"

func (fo *FileOps) register(r *rpn.RPN, shellAvailable bool) {
	r.Register(".", fo.Source, rpn.CatIO, SourceHelp)
	r.Register("append", fo.Append, rpn.CatIO, AppendHelp)
	r.Register("cd", fo.ChangeDir, rpn.CatIO, ChangeDirHelp)
	r.Register("load", fo.Load, rpn.CatIO, LoadHelp)
	r.Register("save", fo.Save, rpn.CatIO, SaveHelp)
	r.Register("source", fo.Source, rpn.CatIO, SourceHelp)
	if shellAvailable {
		r.Register("sh", Shell, rpn.CatIO, ShellHelp)
	}
}
