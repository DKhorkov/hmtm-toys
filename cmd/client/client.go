package main

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/pointers"
	"github.com/DKhorkov/libs/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
)

type Client struct {
	toys.ToysServiceClient
	toys.TagsServiceClient
	toys.MastersServiceClient
	toys.CategoriesServiceClient
}

func main() {
	clientConnection, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", "0.0.0.0", 8060),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		panic(err)
	}

	client := &Client{
		ToysServiceClient:       toys.NewToysServiceClient(clientConnection),
		TagsServiceClient:       toys.NewTagsServiceClient(clientConnection),
		MastersServiceClient:    toys.NewMastersServiceClient(clientConnection),
		CategoriesServiceClient: toys.NewCategoriesServiceClient(clientConnection),
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), requestid.Key, requestid.New())

	tagIDs, err := client.CreateTags(ctx, &toys.CreateTagsIn{
		Tags: []*toys.CreateTagIn{
			{
				Name: "tag1",
			},
			{
				Name: "tag2",
			},
		},
	})
	fmt.Println(tagIDs, err)

	toyByID, err := client.GetToy(ctx, &toys.GetToyIn{ID: 1})
	fmt.Println(err)
	fmt.Println(toyByID)

	allToys, err := client.GetToys(ctx, &toys.GetToysIn{})
	fmt.Println(err)

	for _, toy := range allToys.GetToys() {
		fmt.Println(toy)
	}

	masterID, err := client.RegisterMaster(ctx, &toys.RegisterMasterIn{
		UserID: 30,
		Info:   pointers.New[string]("test master"),
	})
	fmt.Println(err, masterID)

	allMasters, err := client.GetMasters(ctx, &toys.GetMastersIn{})
	fmt.Println(err)

	for _, master := range allMasters.GetMasters() {
		fmt.Println(master)
	}

	toyID, err := client.AddToy(ctx, &toys.AddToyIn{
		UserID:      1,
		Name:        "toy23",
		Price:       120.,
		Quantity:    1,
		CategoryID:  1,
		TagIDs:      []uint32{1},
		Attachments: []string{"someref", "anothererf"},
	})
	fmt.Println(err)
	fmt.Println(toyID)

	master, err := client.GetMasterByUser(ctx, &toys.GetMasterByUserIn{UserID: 1})
	fmt.Println(err)
	fmt.Println(master)

	userToys, err := client.GetUserToys(ctx, &toys.GetUserToysIn{
		UserID: 1,
	})
	fmt.Println("UserToys: ", userToys, err)

	_, err = client.DeleteToy(ctx, &toys.DeleteToyIn{ID: uint64(1)})
	fmt.Println(err)

	_, err = client.UpdateToy(ctx, &toys.UpdateToyIn{
		ID:          2,
		Description: pointers.New[string]("test"),
		Name:        pointers.New[string]("test"),
		CategoryID:  pointers.New[uint32](1),
		Price:       pointers.New[float32](10),
		Quantity:    pointers.New[uint32](1),
		TagIDs:      []uint32{1, 2, 3, 4},
		Attachments: []string{"newRef", "someRef", "anothererf"},
	})
	fmt.Println(err)

	_, err = client.UpdateMaster(ctx, &toys.UpdateMasterIn{
		ID:   masterID.GetMasterID(),
		Info: pointers.New[string]("updated test master info"),
	})
	fmt.Println(err)
}
