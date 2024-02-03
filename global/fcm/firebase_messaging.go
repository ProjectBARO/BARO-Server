package fcm

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var app *firebase.App

func init() {
	fmt.Println("Initializing Firebase App...")

	var err error
	FIREBASE_SERVICE_KEY := os.Getenv("FIREBASE_SERVICE_KEY")
	if FIREBASE_SERVICE_KEY == "" {
		fmt.Println("FIREBASE_SERVICE_KEY is empty")
	}

	opt := option.WithCredentialsJSON([]byte(FIREBASE_SERVICE_KEY))
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
	}
}

func SendPushNotification(fcmToken string, title string, body string) error {
	ctx := context.Background()

	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Data: map[string]string{
			"title": title,
			"body":  body,
		},
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: fcmToken,
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		fmt.Printf("error sending message: %v\n", err)
	}

	return nil
}
