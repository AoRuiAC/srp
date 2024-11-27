package protocol

type contextKey struct {
	name string
}

var ContextKeyReverseProxyAuthed = &contextKey{"rp_authed"}
var ContextKeyProxyAuthed = &contextKey{"p_authed"}

type CachedProxyKey struct {
	Target string
}
