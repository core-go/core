package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

func Connect(ctx context.Context, credentials []byte) (*firestore.Client, error) {
	app, er1 := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(credentials))
	if er1 != nil {
		log.Fatalf("Could not create admin client: %v", er1)
		return nil, er1
	}

	client, er2 := app.Firestore(ctx)
	if er2 != nil {
		log.Fatalf("Could not create data operations client: %v", er2)
		return nil, er2
	}
	return client, nil
}
