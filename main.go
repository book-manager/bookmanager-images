package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/satori/go.uuid"

	"cloud.google.com/go/storage"
)

// Image struct represent incoming requests
type Image struct {
	Image string `json:"image"`
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Post("/author", uploadAuthorAvatar)
	})

	http.ListenAndServe(":3333", r)

}

func uploadAuthorAvatar(w http.ResponseWriter, r *http.Request) {
	var image Image

	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = render.DecodeJSON(r.Body, &image)
	if err != nil {
		log.Fatal(err)
		return
	}
	u1 := uuid.Must(uuid.NewV4())
	fileName := saveToFile(u1.String(), image.Image)
	upload(client, fileName)
}

func upload(client *storage.Client, filename string) error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		fmt.Fprintf(os.Stderr, "BUCKET NAME environment variable must be set.\n")
		os.Exit(1)
	}

	ctx := context.Background()

	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(bucketName).Object(filename).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}

func saveToFile(name string, image string) string {
	elements := strings.Split(image, ",")
	encodedData := elements[len(elements)-1]
	fileName := fmt.Sprintf("%s.jpg", name)

	dec, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	return fileName
}
