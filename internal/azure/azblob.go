package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/logrusorgru/aurora"
)

// colorizer
var (
	style       aurora.Aurora
	accountName string
	accountKey  string
	p           pipeline.Pipeline
)

// auth checks for AZURE_STORAGE_ACCOUNT and AZURE_STORAGE_ACCESS_KEY in the environment
// if env vars are set, it will authenticate and create an azblob request pipeline
// returns false and error message if credentials are not set or are invalid
func auth() error {
	accountName, accountKey = os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) != 0 && len(accountKey) != 0 {
		log.Println("AZURE_STORAGE_ACCOUNT and AZURE_STORAGE_ACCESS_KEY are set in the environment. They will be used to connect to Azure storage.")
		// Create a default request pipeline
		credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err == nil {
			p = azblob.NewPipeline(credential, azblob.PipelineOptions{})
			return nil
		}
		return err

	}
	return fmt.Errorf("either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
}

// ReadFile reads a file from storage container and saves it in a desired location.
func ReadFile(containerName string, filename string, outFile string, noColors bool) {
	style = aurora.NewAurora(!noColors)
	if err := auth(); err != nil {
		log.Fatal(style.Bold(style.Red("ERROR: " + err.Error())))
	}

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)

	ctx := context.Background()

	blobURL := containerURL.NewBlockBlobURL(filename)
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		log.Fatal(style.Bold(style.Red("ERROR: failed to download file " + filename + " with error: " + err.Error())))
	}
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	if _, err = downloadedData.ReadFrom(bodyStream); err != nil {
		log.Fatal(style.Bold(style.Red("ERROR: failed to download file " + filename + " with error: " + err.Error())))
	}

	// create output file and write to it
	var writers []io.Writer
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(style.Bold(style.Red("ERROR: Failed to create an output file: " + err.Error())))
	}
	writers = append(writers, file)
	defer file.Close()

	dest := io.MultiWriter(writers...)
	if _, err := downloadedData.WriteTo(dest); err != nil {
		log.Fatal(style.Bold(style.Red("ERROR: Failed to read object content: " + err.Error())))
	}
	log.Println("INFO: Successfully downloaded " + filename + " from Azure storage as " + outFile)
}
