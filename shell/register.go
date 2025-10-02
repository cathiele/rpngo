package shell

import "mattwach/rpngo/rpn"

func Register(r *rpn.RPN) {
	r.Register(".", Source, rpn.CatIO, SourceHelp)
	r.Register("cd", ChangeDir, rpn.CatIO, ChangeDirHelp)
	r.Register("load", Load, rpn.CatIO, LoadHelp)
	r.Register("save", Save, rpn.CatIO, SaveHelp)
	r.Register("sh", Shell, rpn.CatIO, ShellHelp)
	r.Register("source", Source, rpn.CatIO, SourceHelp)
}
