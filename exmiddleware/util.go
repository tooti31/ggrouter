package exmiddleware

import (
	"github.com/itschip/guildedgo"
	"github.com/tooti31/ggrouter/exrouter"
)

func getChannel(s *guildedgo.Client, channelID string) (*guildedgo.ServerChannel, error) {
	channel, err := s.Channel.GetChannel(channelID)
	if err != nil {
		return s.Channel.GetChannel(channelID)
	}
	return channel, err
}

func getGuild(s *guildedgo.Client, guildID string) (*guildedgo.Server, error) {
	guild, err := s.Server.GetServer(guildID)
	if err != nil {
		return s.Server.GetServer(guildID)
	}
	return guild, err
}

func getMember(s *guildedgo.Client, guildID, userID string) (*guildedgo.ServerMember, error) {
	member, err := s.Members.GetServerMember(guildID, userID)
	if err != nil {
		return s.Members.GetServerMember(guildID, userID)
	}
	return member, err
}

// callCatch calls a catch function with an error
func callCatch(ctx *exrouter.Context, fn CatchFunc, err error) {
	if fn == nil {
		return
	}
	ctx.Set(ctxError, err)
	fn(ctx)
}
