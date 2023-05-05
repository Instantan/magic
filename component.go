package magic

type ComponentFn func(s Socket) AppliedView

func Component(compfn ComponentFn) ComponentFn {
	return func(s Socket) AppliedView {
		sref := s.clone()
		res := compfn(sref)
		return res
	}
}
