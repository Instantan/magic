package magic

type ComponentFn[Props any] func(s Socket, p Props) AppliedView

func Component[Props any](compfn ComponentFn[Props]) ComponentFn[Props] {
	return func(s Socket, p Props) AppliedView {
		sref := s.clone()
		res := compfn(sref, p)
		return res
	}
}
