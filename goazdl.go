package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Global          GlobalConfig   `yaml:"global"`
	StorageAccounts []StorageEntry `yaml:"storage_accounts"`
}

type GlobalConfig struct {
	LogfileDirectory string   `yaml:"logfile_directory"`
	SendgridAPIKey   string   `yaml:"sendgrid_api_key"`
	ToMail           []string `yaml:"to_mail"`
}

type StorageEntry struct {
	StorageAccountName string `yaml:"storage_account_name"`
	TenantID           string `yaml:"tenant_id"`
	ClientID           string `yaml:"client_id"`
	ClientSecret       string `yaml:"client_secret"`
	DLPath             string `yaml:"dl_path"`
}

func listBlobs(accountName, accountKey, containerName, dlPath string) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatalf("Invalid Credentials: %v", err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceURL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	containerURL := azblob.NewContainerURL(*serviceURL, p)

	ctx := context.Background()
	fmt.Printf("Blobs in container '%s':\n", containerName)
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, blob := range listBlob.Segment.BlobItems {
			fmt.Println(blob.Name)
		}
		marker = listBlob.NextMarker
	}
}

func main() {
	fmt.Println("Hello, World!")

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	fmt.Printf("Logfile Storage Directory: %s\n", config.Global.LogfileDirectory)
	fmt.Printf("Sendgrid API Key: %s\n", config.Global.SendgridAPIKey)
	fmt.Printf("To Mail Addresses: %s\n", config.Global.ToMail)

	for idx, storage := range config.StorageAccounts {
		fmt.Printf("\nStorage Account %d:\n", idx+1)
		fmt.Println("  Account Name:", storage.StorageAccountName)
		fmt.Println("  Tenant ID:", storage.TenantID)
		fmt.Println("  Client ID:", storage.ClientID)
		fmt.Println("  Client Secret:", storage.ClientSecret)
		fmt.Println("  Download Path:", storage.DLPath)

		listBlobs(
			storage.StorageAccountName,
			storage.ClientSecret,
			storage.StorageAccountName,
			storage.DLPath,
		)
	}
}
