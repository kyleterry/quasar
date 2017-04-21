package quasar

func PrivmsgMiddlware(handler Handler) Handler {
	return HandlerFunc(func(r Result, msg Message, com Communication) {
		if msg.Command == "PRIVMSG" {
			handler.HandleMatch(r, msg, com)
		}
	})
}
