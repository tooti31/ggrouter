package exmiddleware

import (
	"github.com/itschip/guildedgo"
	"github.com/tooti31/ggrouter/exrouter"
)

// Err retrieves the error variable from the context
func Err(ctx *exrouter.Context) error {
	if v := ctx.Get(ctxError); v != nil {
		return v.(error)
	}
	return nil
}

// Guild retrieves the guild variable from a context
func Guild(ctx *exrouter.Context) *guildedgo.Server {
	if v := ctx.Get(ctxGuild); v != nil {
		return v.(*guildedgo.Server)
	}
	return nil
}

// Channel retrieves the channel variable from a context
func Channel(ctx *exrouter.Context) *guildedgo.ServerChannel {
	if v := ctx.Get(ctxChannel); v != nil {
		return v.(*guildedgo.ServerChannel)
	}
	return nil
}

// Member fetches the member from the context
func Member(ctx *exrouter.Context) *guildedgo.ServerMember {
	if v := ctx.Get(ctxMember); v != nil {
		return v.(*guildedgo.ServerMember)
	}
	return nil
}
