package store

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// AzureBlobStorage represents an Azure Storage object.
type AzureBlobStorage struct {
	container string
	object    string
}

// NewAzureBlobStorage returns an Azure storage item given a container and an object
// path.
func NewAzureBlobStorage(container string, object string) *AzureBlobStorage {
	return &AzureBlobStorage{
		container: container,
		object:    object,
	}
}

// Load downloads an object from a bucket and returns an error if it cannot
// be read.
func (s *AzureBlobStorage) Load() ([]byte, error) {
	// Get container URL.
	containerURL, err := s.containerURL()
	if err != nil {
		return nil, err
	}

	// Upload a blob.
	blobURL := containerURL.NewBlockBlobURL(s.object)

	// Download the blob.
	ctx := context.Background()
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	// Automatically retries are performed if the connection fails.
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	// Read from buffer.
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return nil, err
	}

	return downloadedData.Bytes(), nil
}

// Save uploads an object to a bucket and returns an error if it cannot be
// written.
func (s *AzureBlobStorage) Save(b []byte) error {
	// Get container URL.
	containerURL, err := s.containerURL()
	if err != nil {
		return err
	}

	// Upload a blob.
	blobURL := containerURL.NewBlockBlobURL(s.object)

	// Upload file.
	ctx := context.Background()
	_, err = blobURL.Upload(ctx, bytes.NewReader(b), azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{}, azblob.AccessTierHot, azblob.BlobTagsMap{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (s *AzureBlobStorage) containerURL() (azblob.ContainerURL, error) {
	accountName := ""
	accountKey := ""

	// Get the Azure credentials.
	connString := os.Getenv("AzureWebJobsStorage")
	if len(connString) == 0 {
		// Get the storage account name and key and set environment variables.
		accountName, accountKey = os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
		if len(accountName) == 0 || len(accountKey) == 0 {
			return azblob.ContainerURL{}, fmt.Errorf("either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
		}
	} else {
		// Parse the connection string.
		arr := strings.Split(connString, ";")
		for _, v := range arr {
			pair := strings.SplitN(v, "=", 2)
			if len(pair) < 2 {
				continue
			}

			switch pair[0] {
			case "AccountName":
				accountName = pair[1]
			case "AccountKey":
				accountKey = pair[1]
			}
		}
	}

	// Create a default request pipeline using the storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return azblob.ContainerURL{}, err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Build the storage account blob service URL endpoint.
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, s.container))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)

	return containerURL, nil
}
