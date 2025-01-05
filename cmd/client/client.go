package main

import (
	"context"
	"fmt"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/libs/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	requestID := requestid.New()

	// toyByID, err := client.GetToy(context.Background(), &toys.GetToyIn{ID: 1})
	// fmt.Println(err)
	// fmt.Println(toyByID)

	allToys, err := client.GetToys(context.Background(), &toys.GetToysIn{RequestID: requestID})
	fmt.Println(err)
	for _, toy := range allToys.GetToys() {
		fmt.Println(toy)
	}

	// masterID, err := client.RegisterMaster(context.Background(), &toys.RegisterMasterIn{
	//	UserID: 1,
	//	Info:   "testInfo",
	// })
	// fmt.Println(err, masterID)

	// allMasters, err := client.GetMasters(context.Background(), &toys.GetMastersIn{RequestID: requestID})
	// fmt.Println(err)
	// for _, master := range allMasters.GetMasters() {
	//	fmt.Println(master)
	//}

	// toyID, err := client.AddToy(context.Background(), &toys.AddToyIn{
	//	RequestID:   requestID,
	//	UserID:      1,
	//	Name:        "toy23",
	//	Price:       120.,
	//	Quantity:    1,
	//	CategoryID:  1,
	//	TagIDs:      []uint32{1},
	//	Attachments: []string{"someref", "anothererf"},
	// })
	// fmt.Println(err)
	// fmt.Println(toyID)

	// master, err := client.GetMasterByUser(context.Background(), &toys.GetMasterByUserIn{UserID: 1})
	// fmt.Println(err)
	// fmt.Println(master)

	// userToys, err := client.GetUserToys(context.Background(), &toys.GetUserToysIn{
	//	RequestID: requestID,
	//	UserID:    4,
	// })
	// fmt.Println("UserToys: ", userToys, err)
}
