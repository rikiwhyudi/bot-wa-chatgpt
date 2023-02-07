package handler

import (
	"context"
	"fmt"
	"os"
	"strings"

	gogpt "github.com/sashabaranov/go-gpt3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func GetEventHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			if !v.Info.IsFromMe && !v.Info.IsGroup {
				if v.Message.GetConversation() != "" {
					var messageBody = v.Message.GetConversation()
					fmt.Println("Received a message!", messageBody)

					response, err := GenerateResponse(messageBody)
					if err != nil {
						response = "I'm sorry, an error occurred while processing your request."
					}

					client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(strings.TrimSpace(response)),
					})
				}
			}
		}
	}
}

func GenerateResponse(message string) (string, error) {
	c := gogpt.NewClient(os.Getenv("API_KEY"))
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           message,
		Temperature:      0.01,
		MaxTokens:        450,
		TopP:             1,
		FrequencyPenalty: 0.5,
		PresencePenalty:  0.06,
	}

	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Text, nil

}
