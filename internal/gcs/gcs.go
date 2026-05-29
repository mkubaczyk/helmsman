package gcs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	// Imports the Google Cloud Storage client package.
	"cloud.google.com/go/storage"
	"github.com/logrusorgru/aurora"
	netContext "golang.org/x/net/context"
)

// colorizer
var style aurora.Aurora

func IsRunningInGCP() bool {
	resp, err := http.Get("http://metadata.google.internal")

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	return err == nil
}

// Auth checks for GCLOUD_CREDENTIALS in the environment
// returns true if they exist and creates a json credentials file and sets the GOOGLE_APPLICATION_CREDENTIALS env var
// returns true if GCP metadata server is present
// returns false if credentials are not found
func Auth() (string, error) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		return "GOOGLE_APPLICATION_CREDENTIALS is already set in the environment", nil
	}

	if os.Getenv("GCLOUD_CREDENTIALS") != "" {
		credFile := "/tmp/gcloud_credentials.json"
		// write the credentials content into a json file
		d := []byte(os.Getenv("GCLOUD_CREDENTIALS"))
		err := ioutil.WriteFile(credFile, d, 0o644)
		if err != nil {
			return fmt.Sprintf("Cannot create credentials file: %s", err), err
		}

		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credFile)
		return "ok", nil
	}

	if IsRunningInGCP() {
		return "Metadata server present, running in GCP", nil
	}

	return "can't authenticate", fmt.Errorf("can't authenticate")
}

// ReadFile reads a file from storage bucket and saves it in a desired location.
func ReadFile(bucketName string, filename string, outFile string, noColors bool) (string, error) {
	style = aurora.NewAurora(!noColors)
	if msg, err := Auth(); err != nil {
		return msg, nil
	}

	ctx := netContext.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "Failed to configure Storage bucket: ", err
	}
	storageBucket := client.Bucket(bucketName)

	// Creates an Object handler for our file
	obj := storageBucket.Object(filename)

	// Read the object.
	r, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Sprintf("Failed to create object reader: %s", err), err
	}
	defer r.Close()

	// create output file and write to it
	var writers []io.Writer
	file, err := os.Create(outFile)
	if err != nil {
		return fmt.Sprintf("Failed to create an output file: %s", err), err
	}
	writers = append(writers, file)
	defer file.Close()

	dest := io.MultiWriter(writers...)
	if _, err := io.Copy(dest, r); err != nil {
		return fmt.Sprintf("Failed to read object content: %s", err), err
	}
	return "Successfully downloaded " + filename + " from GCS as " + outFile, nil
}
