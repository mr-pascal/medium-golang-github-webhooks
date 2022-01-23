package main

import (
	"github.com/go-playground/webhooks/v6/github"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"log"
	"net/http"
)

func main() {
	// Initialize GitHub webhook parser
	hook, _ := github.New(github.Options.Secret("my_secret"))

	// Create new Fiber application
	app := fiber.New()

	// Post endpoint for handling webhooks
	app.Post("/webhook", func(c *fiber.Ctx) error {
		log.Println("Received webhook...")

		// Convert fasthttp.Request to http.Request
		r := new(http.Request)
		fasthttpadaptor.ConvertRequest(c.Context(), r, true)

		// Parse request and extract event payload
		payload, e := hook.Parse(r, github.IssueCommentEvent)
		if e != nil {
			log.Println("Error parsing", e)
		}

		switch payload.(type) {
		// Handling new comment on issue event
		case github.IssueCommentPayload:
			comment := payload.(github.IssueCommentPayload)

			// Retrieve individual fields of interest
			commentText := comment.Comment.Body
			userName := comment.Comment.User.Login

			// Print the comment creator and the comment text
			log.Printf("User '%s' posted '%s'", userName, commentText)
		}

		return c.SendStatus(200)
	})

	// Listen on the port '3000'
	app.Listen(":3000")
}
