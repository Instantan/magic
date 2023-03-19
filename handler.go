package magic

type StaticHandler interface {
}

type LiveHandler interface {
	OnMount()
	OnEvent()
	OnUnmount()
}

func HandleStatic(handler StaticHandler) {

}

func HandleLive(handler LiveHandler) {

}
