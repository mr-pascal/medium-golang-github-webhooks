package main

/**
Import needed modules
*/
import (
	"github.com/go-playground/webhooks/v6/github"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"log"
	"net/http"
)

func main() {

	/**
	### Initialize GitHub webhook parser ###
	-   This parser helps takes care about parsing the plain HTTP(S)
	    request to the GitHub event struct the application needs to work with.
	-   Also, it (optionally) can check for a secret that can
	    be defined for webhooks.
	-   In this case the parser will also verify that the incoming request contains the secret "my_secret"
	    otherwise it will error.
	*/
	hook, _ := github.New(github.Options.Secret("my_secret"))

	/**
	Create new Go Fiber application
	*/
	app := fiber.New()

	/**
	Create POST endpoint on the "/webhook" route to handle incoming webhooks
	*/
	app.Post("/webhook", func(c *fiber.Ctx) error {
		log.Println("Received webhook...")

		/**
		### Convert fasthttp.Request to http.Request ###
		-   The request object needs to be converted, because the GitHub
		    webhook parser explicitly needs an "*http.Request" as input parameter.
		*/
		httpRequest := new(http.Request)
		err := fasthttpadaptor.ConvertRequest(c.Context(), httpRequest, true)
		if err != nil {
			log.Println("Error converting request", err)
		}

		/**
		### Verify and parse the defined events inside the request ###
		-   The `.Parse(...)` method checks the request and try to parse
		    the defined events. In this case only one event `github.IssueCommentEvent`
		    was defined, since this is the only one we are interested in.
		-   In this step it is also checked if the above defined `my_secret` is
		    set correctly in the webhook request.
		*/
		payload, e := hook.Parse(httpRequest, github.IssueCommentEvent)
		if e != nil {
			log.Println("Error parsing", e)
		}

		/**
		Switch case to apply business logic based on the event received
		via the webhook.
		*/
		switch payload.(type) {
		/**
		Handling "Comment on Issue" event.
		*/
		case github.IssueCommentPayload:
			/**
			Casting `payload` into the `github.IssueCommentPayload` struct
			*/
			comment := payload.(github.IssueCommentPayload)

			/**
			### Retrieve individual fields of interest ###
			-   The actual `github.IssueCommentPayload` struct is very big. But
			    we are only interested in the content the user posted and their name.
			*/
			commentText := comment.Comment.Body
			userName := comment.Comment.User.Login

			/**
			Print the comment creator and the posted text to the terminal
			*/
			log.Printf("User '%s' posted '%s'", userName, commentText)
		}

		/**
		### Respond to the webhook ###
		-   In this case we are always acknowledging the webhook with a `200`.
		-   In a real case you might want to return a `4XX` or `5XX` depending
		    on your actual business logic.
		*/
		return c.SendStatus(200)
	})

	/**
	Listen on the port `3000`
	*/
	app.Listen(":3000")
}
