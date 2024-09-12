package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gocarina/gocsv"
	"github.com/subosito/gotenv"
	"google.golang.org/api/option"
)

func init() {
	gotenv.Load()
}

func main() {
	GetFileWithGeneration()
}

// GetFileWithGeneration function
func GetFileWithGeneration() error {
	bucket := {BUCKET_NAME}"
	var gen int64 = {GENERATION_NUMBER}
	fileName := "{FILE_NAME}.csv"

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	fmt.Printf("File : %v, Generation : %v\n", fileName, gen)
	obj := client.Bucket(bucket).Object(fileName)
	// set file generation
	obj = obj.Generation(gen)

	storeReader, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	var rec []Song
	reader := csv.NewReader(storeReader)
	if err := gocsv.UnmarshalCSV(reader, &rec); err != nil {
		panic(err)
	}

	defer storeReader.Close()

	dataByte, err := json.Marshal(rec)
	fmt.Printf("Print data from csv = %v\n", string(dataByte))

	return nil
}

type Song struct {
	ID            string `csv:"id" json:"id"`
	Name          string `csv:"name" json:"name"`
	Size          string `csv:"size" json:"size"`
	Category      string `csv:"category" json:"category"`
}
