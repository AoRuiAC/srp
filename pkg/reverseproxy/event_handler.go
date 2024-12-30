package reverseproxy

type EventHandler struct {
	OnAdd    func(host string, port string)
	OnRemove func(host string, port string)
}

type EventHandlers []EventHandler

func (hs EventHandlers) OnAdd(host string, port string) {
	for _, h := range hs {
		if h.OnAdd != nil {
			go h.OnAdd(host, port)
		}
	}
}

func (hs EventHandlers) OnRemove(host string, port string) {
	for _, h := range hs {
		if h.OnRemove != nil {
			go h.OnRemove(host, port)
		}
	}
}
