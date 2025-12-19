package aws

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var _ = flag.String("f", "", "") // Accept -f flag from Makefile

// Verify AWS SDK returns an error when no valid credentials are found.
func TestS3DownloadFailsWithoutCredentials(t *testing.T) {
	// Clear AWS-related env vars and restore them later
	envVars := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
		"AWS_PROFILE",
		"AWS_SHARED_CREDENTIALS_FILE",
		"AWS_CONFIG_FILE",
	}
	saved := make(map[string]string)
	for _, v := range envVars {
		saved[v] = os.Getenv(v)
		os.Unsetenv(v)
	}
	defer func() {
		for k, v := range saved {
			if v != "" {
				os.Setenv(k, v)
			}
		}
	}()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewCredentials(&credentials.ChainProvider{
			Providers: []credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{Filename: "/nonexistent", Profile: "nonexistent"},
			},
			VerboseErrors: true,
		}),
	})
	if err != nil {
		t.Fatalf("session creation failed: %v", err)
	}

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(&fakeWriterAt{}, &s3.GetObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-key"),
	})

	if err == nil {
		t.Fatal("expected error when no credentials available, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "NoCredentialProviders") {
		t.Errorf("expected NoCredentialProviders error, got: %v", err)
	}
}

type fakeWriterAt struct{}

func (f *fakeWriterAt) WriteAt(p []byte, off int64) (n int, err error) {
	return len(p), nil
}
