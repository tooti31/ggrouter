package exrouter

import (
	"fmt"
	"github.com/tooti31/ggrouter"
	"sync"

	"github.com/itschip/guildedgo"
)

// Context represents a command context
type Context struct {
	// Route is the route that this command came from
	Route *dgrouter.Route
	Msg   *guildedgo.ChatMessage
	Ses   *guildedgo.Client

	// List of arguments supplied with the command
	Args Args

	// Vars that can be optionally set using the Set and Get functions
	vmu  sync.RWMutex
	Vars map[string]interface{}
}

// Set sets a variable on the context
func (c *Context) Set(key string, d interface{}) {
	c.vmu.Lock()
	c.Vars[key] = d
	c.vmu.Unlock()
}

// Get retrieves a variable from the context
func (c *Context) Get(key string) interface{} {
	if c, ok := c.Vars[key]; ok {
		return c
	}
	return nil
}

// Reply replies to the sender with the given message
func (c *Context) Reply(args ...interface{}) (*guildedgo.ChatMessage, error) {
	return c.Ses.Channel.SendMessage(c.Msg.ChannelID, &guildedgo.MessageObject{Content: fmt.Sprint(args...), ReplyMessageIds: []string{c.Msg.ID}})
}

// ReplyEmbed replies to the sender with an embed
func (c *Context) ReplyEmbed(args ...interface{}) (*guildedgo.ChatMessage, error) {
	return c.Ses.Channel.SendMessage(c.Msg.ChannelID, &guildedgo.MessageObject{
		ReplyMessageIds: []string{c.Msg.ID},
		Embeds: []guildedgo.ChatEmbed{{
			Description: fmt.Sprint(args...),
		}},
	})
}

// Guild retrieves a guild from the state or restapi
func (c *Context) Guild(guildID string) (*guildedgo.Server, error) {
	g, err := c.Ses.Server.GetServer(guildID)
	if err != nil {
		g, err = c.Guild(guildID)
	}
	return g, err
}

// Channel retrieves a channel from the state or restapi
func (c *Context) Channel(channelID string) (*guildedgo.ServerChannel, error) {
	g, err := c.Ses.Channel.GetChannel(channelID)
	if err != nil {
		g, err = c.Channel(channelID)
	}
	return g, err
}

// Member retrieves a member from the state or restapi
func (c *Context) Member(guildID, userID string) (*guildedgo.User, error) {
	m, err := c.Member(guildID, userID)
	if err != nil {
		m, err = c.Member(guildID, userID)
	}
	return m, err
}

// NewContext returns a new context from a message
func NewContext(s *guildedgo.Client, m *guildedgo.ChatMessage, args Args, route *dgrouter.Route) *Context {
	return &Context{
		Route: route,
		Msg:   m,
		Ses:   s,
		Args:  args,
		Vars:  map[string]interface{}{},
	}
}
