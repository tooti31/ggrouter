package exrouter

import (
	"github.com/itschip/guildedgo"
	"github.com/tooti31/ggrouter"
	"strings"
)

// HandlerFunc ...
type HandlerFunc func(*Context)

// MiddlewareFunc is middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Route wraps dgrouter.Router to use a Context
type Route struct {
	*dgrouter.Route
}

// New returns a new router wrapper
func New() *Route {
	return &Route{
		Route: dgrouter.New(),
	}
}

// On registers a handler function
func (r *Route) On(name string, handler HandlerFunc) *Route {
	return &Route{r.Route.On(name, WrapHandler(handler))}
}

// Group ...
func (r *Route) Group(fn func(rt *Route)) *Route {
	return &Route{r.Route.Group(func(r *dgrouter.Route) {
		fn(&Route{r})
	})}
}

// Use ...
func (r *Route) Use(fn ...MiddlewareFunc) *Route {
	wrapped := []dgrouter.MiddlewareFunc{}
	for _, v := range fn {
		wrapped = append(wrapped, WrapMiddleware(v))
	}
	return &Route{
		r.Route.Use(wrapped...),
	}
}

// WrapMiddleware ...
func WrapMiddleware(mware MiddlewareFunc) dgrouter.MiddlewareFunc {
	return func(next dgrouter.HandlerFunc) dgrouter.HandlerFunc {
		return func(i interface{}) {
			WrapHandler(mware(UnwrapHandler(next)))(i)
		}
	}
}

// OnMatch registers a route with the given matcher
func (r *Route) OnMatch(name string, matcher func(string) bool, handler HandlerFunc) *Route {
	return &Route{r.Route.OnMatch(name, matcher, WrapHandler(handler))}
}

func mention(s *guildedgo.Client, ids []guildedgo.MentionsUser) string {
	user, _ := s.Users.GetOwnUser()
	for _, v := range ids {
		if user.Id == v.ID {
			return user.Id
		}
	}
	return ""
}

// FindAndExecute is a helper method for calling routes
// it creates a context from a message, finds its route, and executes the handler
// it looks for a message prefix which is either the prefix specified or the message is prefixed
// with a bot mention
//
//	s            : guildedgo Client to pass to context
//	prefix       : prefix you want the bot to respond to
//	botID        : user ID of the bot to allow you to substitute the bot ID for a prefix
//	m            : guildedgo.ChatMessage to pass to context
func (r *Route) FindAndExecute(s *guildedgo.Client, prefix string, m *guildedgo.ChatMessage) error {
	var pf string

	p := func(t string) bool {
		return strings.HasPrefix(m.Content, t)
	}

	switch {
	case prefix != "" && p(prefix):
		pf = prefix
	case p(mention(s, m.Mentions.Users)):
		pf = mention(s, m.Mentions.Users)
	default:
		return dgrouter.ErrCouldNotFindRoute
	}

	command := strings.TrimPrefix(m.Content, pf)
	args := ParseArgs(command)

	if rt, depth := r.FindFull(args...); depth > 0 {
		args = append([]string{strings.Join(args[:depth], string(separator))}, args[depth:]...)
		rt.Handler(NewContext(s, m, args, rt))
	} else {
		return dgrouter.ErrCouldNotFindRoute
	}

	return nil
}

// WrapHandler wraps a dgrouter.HandlerFunc
func WrapHandler(fn HandlerFunc) dgrouter.HandlerFunc {
	if fn == nil {
		return nil
	}
	return func(i interface{}) {
		fn(i.(*Context))
	}
}

// UnwrapHandler unwraps a handler
func UnwrapHandler(fn dgrouter.HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		fn(ctx)
	}
}
