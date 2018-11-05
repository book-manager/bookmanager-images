# BookManager Images

Upload incoming images to Google storage

# Install

```
go get github.com/go-chi/chi
go get github.com/go-chi/render
go get github.com/satori/go.uuid
go get cloud.google.com/go/storage
```

# Build

`go build -o main`

# Run

You need to have correct credentials in order to upload files to Google Storage.

## Environment variables

`GOOGLE_APPLICATION_CREDENTIALS` path to a JSON file with credentials

`GOOGLE_CLOUD_PROJECT` name of a Google project in which there is a bucket

`BUCKET_NAME` name of the bucket used to upload images


# Docker

`docker run -it --rm -p 3333:3333 -e GOOGLE_APPLICATION_CREDENTIALS="path-to-credentials" -e GOOGLE_CLOUD_PROJECT="project-name" -e BUCKET_NAME="bucket-name" cosaquee/bookmanager-images`

Docker volumes is the best option to provide GCP credentials.