package reverseproxy

type EventHandler struct {
	OnAdd    func(string)
	OnRemove func(string)
}

type EventHandlers []EventHandler

func (hs EventHandlers) OnAdd(socket string) {
	for _, h := range hs {
		if h.OnAdd != nil {
			go h.OnAdd(socket)
		}
	}
}

func (hs EventHandlers) OnRemove(socket string) {
	for _, h := range hs {
		if h.OnRemove != nil {
			go h.OnRemove(socket)
		}
	}
}
