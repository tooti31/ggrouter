package exrouter_test

import (
	dgrouter "github.com/tooti31/ggrouter"
	"github.com/tooti31/ggrouter/exrouter"
	"log"
	"testing"

	"github.com/itschip/guildedgo"
)

func TestRouter(t *testing.T) {
	messages := []string{
		"!ping",
		"!say hello",
		"!test args one two three",
	}

	r := exrouter.Route{
		Route: dgrouter.New(),
	}

	r.On("ping", func(ctx *exrouter.Context) {})

	r.On("say", func(ctx *exrouter.Context) {
		if ctx.Args.Get(1) != "hello" {
			t.Log("say fail")
			t.Fail()
		}
	})

	r.On("test", func(ctx *exrouter.Context) {
		ctx.Set("hello", "hi")
		if r := ctx.Get("hello"); r.(string) != "hi" {
			t.Log("test fail")
			t.Fail()
		}
		expected := []string{"args", "one", "two", "three"}
		for i, v := range expected {
			if ctx.Args.Get(i+1) != v {
				t.Log("args fail")
				t.Fail()
			}
		}
	})

	r.On("help", func(ctx *exrouter.Context) {
		log.Println("Bot was mentioned")
	})

	for _, v := range messages {
		// Construct mock message
		msg := &guildedgo.ChatMessage{
			ID:        "xcwcqwc",
			Type:      "bot",
			ServerID:  "greg3",
			ChannelID: "f32f32",
			Content:   v,
		}

		// Attempt to find and execute the route for this message
		err := r.FindAndExecute(nil, "!", msg)
		if err != nil {
			t.Log("FindAndExecute fail")
			t.Fail()
		}
	}
}
