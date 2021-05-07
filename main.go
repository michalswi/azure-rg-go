package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type clientInfo struct {
	SubscriptionID string
}

type RG struct {
	Name    string   `json:"name"`
	Details RGdetais `json:"details"`
}

type RGdetais struct {
	ID                string `json:"id"`
	Location          string `json:"location"`
	ProvisioningState string `json:"state"`
	Tags              string `json:"tags"`
}

var (
	ctx                   = context.Background()
	clientData            clientInfo
	authorizer            autorest.Authorizer
	subs                  = os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroupName     = "goRGtest"
	resourceGroupLocation = "westeurope"
)

func main() {
	clientData.SubscriptionID = subs
	var err error
	// https://github.com/MicrosoftDocs/azure-dev-docs/issues/496
	authorizer, err = auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Fatalf("AuthN failed: %v", err)
	}
	if len(os.Args) < 2 {
		log.Fatalf("Missing args. Available create/list")
	} else {
		switch os.Args[1] {
		case "create":
			groupCreate, err := createRGroup()
			if err != nil {
				log.Fatalf("Failed to create group: %v", err)
			}
			log.Printf("Created resource group: %v", *groupCreate.Name)
		case "list":
			listRGroup()
		case "delete":
			deleteRGroup()
		}
	}
}

func createRGroup() (group resources.Group, err error) {
	groupsClient := resources.NewGroupsClient(clientData.SubscriptionID)
	groupsClient.Authorizer = authorizer
	return groupsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		resources.Group{
			Location: to.StringPtr(resourceGroupLocation)})
}

func deleteRGroup() {
	groupsClient := resources.NewGroupsClient(clientData.SubscriptionID)
	groupsClient.Authorizer = authorizer
	_, err := groupsClient.Delete(ctx, resourceGroupName)
	if err != nil {
		log.Fatalf("Delete RG failed: %v", err)
	}
	log.Printf("Deleted resource group: %s", resourceGroupName)
}

func listRGroup() {
	env, _ := azure.EnvironmentFromName("AzurePublicCloud")
	resourceGroupsClient := resources.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, subs)
	resourceGroupsClient.Authorizer = authorizer
	groupsList, err := resourceGroupsClient.List(ctx, "", nil)
	if err != nil {
		log.Fatalf("List RG failed: %v", err)
	}
	for _, g := range groupsList.Values() {
		t := ""
		if g.Tags == nil || len(g.Tags) <= 0 {
			t += "no tags"
		} else {
			for k, v := range g.Tags {
				t += fmt.Sprintf("\t%s = %s\n", k, *v)
			}
		}
		elements := &RG{
			Name: *g.Name,
			Details: RGdetais{
				ID:                *g.ID,
				Location:          *g.Location,
				ProvisioningState: *g.Properties.ProvisioningState,
				Tags:              t,
			},
		}
		b, _ := json.Marshal(elements)
		fmt.Printf("%s", b)
	}
}
