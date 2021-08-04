package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func main() {
	subscriptionID := os.Args[1]
	AzureAuth(subscriptionID)
}

func AzureAuth(subscriptionID string) compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Auth: Successful")
		vmClient.Authorizer = authorizer
	}
	return vmClient
}
