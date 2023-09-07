package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/minio/minio-go/v6"

	"toktik/pkg/config"
)

var configPath string

func init() {
	go runMinio(os.Args[3:])

	flag.StringVar(&configPath, "config", "/etc/config.yaml", "config path")
}

func main() {
	flag.Parse()
	config.ReadConfigFromLocal(configPath)

	var client *minio.Client
	var err error
	// Connect to minio
	client, err = minio.New(
		config.Conf.GetString(config.KEY_MINIO_ENDPOINT),
		config.Conf.GetString(config.KEY_MINIO_ACCESS_KEY),
		config.Conf.GetString(config.KEY_MINIO_SECRET_KEY),
		false,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Create buckets
	buckets := config.Conf.Get("minio.buckets").([]any)
	bucketList := make([]string, 0, len(buckets))
	location := config.Conf.GetString("minio.location")
	for _, bucket := range buckets {
		bucketList = append(bucketList, bucket.(string))
	}

	retryTime := 0
	for {
		err := createBuckets(client, bucketList, location)
		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "connection refused") {
				retryTime++
				log.Printf("create buckets timeout, retry %d times\n", retryTime)
			} else {
				log.Fatalln(err)
			}
		} else {
			break
		}
	}

	log.Println("Init Minio success")

	select {}

	// Connect to rabbitmq
}

func createBuckets(client *minio.Client, buckets []string, location string) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	for _, bucket := range buckets {
		exist, err := client.BucketExistsWithContext(ctx, bucket)
		if err != nil {
			<-ctx.Done()
			return err
		}
		if !exist {
			err := client.MakeBucketWithContext(ctx, bucket, location)
			if err != nil {
				return err
			}
			err = client.SetBucketPolicy(bucket, `{"Version":"2012-10-17","Statement":[{"Action":["s3:GetBucketLocation","s3:ListBucket"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::`+bucket+`"],"Sid":""},{"Action":["s3:GetObject"],"Effect":"Allow","Principal":{"AWS":["*"]},"Resource":["arn:aws:s3:::`+bucket+`/*"],"Sid":""}]}`)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runMinio(args []string) {
	cmd := exec.Command("minio", args...)
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
