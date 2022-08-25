package try

type TryCatch struct {
	handler        func()
	catchHandler   func(interface{})
	finallyHandler func()
}

func NewTry(handler func()) *TryCatch {
	tryObj := new(TryCatch)
	tryObj.handler = handler
	return tryObj
}

func (ptl *TryCatch) Catch(handler func(err interface{})) *TryCatch {
	ptl.catchHandler = handler
	return ptl
}

func (ptl *TryCatch) Finally(handler func()) *TryCatch {
	ptl.finallyHandler = handler
	return ptl
}

func (ptl *TryCatch) Run() {
	defer func() {
		if err := recover(); err != nil {
			if ptl.catchHandler != nil {
				ptl.catchHandler(err)
			}
			if ptl.finallyHandler != nil {
				ptl.finallyHandler()
			}
		}
	}()
	ptl.handler()

}
