package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

func InitCosmos() func() {
	f, err := os.Create("cosmos-log-file.txt")
	handle(err)

	// Configure the listener to write to a file rather than to the console
	azlog.SetListener(func(event azlog.Event, s string) {
		f.WriteString(s + "\n")
	})

	// Filter the types of events you'd like to log by removing the ones you're not interested in (if any)
	// We recommend using the default logging with no filters - but if filtering we recommend *always* including
	// `azlog.EventResponseError` since this is the event type that will help with debugging errors
	azlog.SetEvents(azlog.EventRequest, azlog.EventResponse, azlog.EventRetryPolicy, azlog.EventResponseError)

	return func() {
		f.Close()
	}
}

func RunCosmos() {
	ctx := context.TODO()

	// Using account keys
	var (
		cosmosDbEndpoint = os.Getenv("COSMOS_ENDPOINT")
		cosmosDbKey      = os.Getenv("COSMOS_KEY")
	)

	cred, err := azcosmos.NewKeyCredential(cosmosDbKey)
	handle(err)
	client, err := azcosmos.NewClientWithKey(cosmosDbEndpoint, cred, nil)
	handle(err)

	tenantId := "dev"
	id, err := uuid.NewV7()
	handle(err)

	// CRUD operation on Items
	item := map[string]string{
		"tenantId": tenantId,
		"id":       id.String(),
		"value":    "2",
	}

	marshalled, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}

	container, err := client.NewContainer(os.Getenv("COSMOS_DB"), os.Getenv("COSMOS_CONTAINER"))
	handle(err)

	pk := azcosmos.NewPartitionKeyString(tenantId)

	// Create an item
	itemResponse, err := container.CreateItem(ctx, pk, marshalled, nil)
	handle(err)

	// Read an item
	itemResponse, err = container.ReadItem(ctx, pk, id.String(), nil)
	handle(err)

	var itemResponseBody map[string]string
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		log.Print(err)
	}

	itemResponseBody["value"] = "3"
	marshalledReplace, err := json.Marshal(itemResponseBody)
	if err != nil {
		log.Print(err)
	}

	// Replace an item
	itemResponse, err = container.ReplaceItem(ctx, pk, id.String(), marshalledReplace, nil)
	handle(err)

	// Patch an item
	patch := azcosmos.PatchOperations{}
	patch.AppendAdd("/newField", "newValue")
	patch.AppendRemove("/oldFieldToRemove")

	itemResponse, err = container.PatchItem(ctx, pk, id.String(), patch, nil)
	handle(err)

	// Delete an item
	itemResponse, err = container.DeleteItem(ctx, pk, id.String(), nil)
	handle(err)
}
