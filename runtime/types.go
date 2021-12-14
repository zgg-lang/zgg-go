package runtime

type IEval interface {
	Eval(*Context)
}
