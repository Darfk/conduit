package conduit

// generated from divide(in *divop) (out *divop)
func divideStage(inc <-chan *divop, cancel <-chan struct{}) <-chan *divop {
	ouc := make(chan *divop)
	go func() {
		defer close(ouc)
		for in := range inc {
			ouv := divide(in)
			select {
			case <-cancel:
				return
			case ouc <- ouv:
			}
		}
	}()
	return ouc
}
