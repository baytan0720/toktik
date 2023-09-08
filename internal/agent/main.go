package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/minio/minio-go/v6"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	"toktik/pkg/config"
	"toktik/pkg/rabbitmq"
)

var configPath string

func init() {
	go runMinio(os.Args[3:])

	flag.StringVar(&configPath, "config", "/etc/config.yaml", "config path")
}

func main() {
	flag.Parse()
	config.ReadConfigFromLocal(configPath)

	// Connect to minio
	client, err := minio.New(
		config.GetString(config.KEY_MINIO_ENDPOINT),
		config.GetString(config.KEY_MINIO_ACCESS_KEY),
		config.GetString(config.KEY_MINIO_SECRET_KEY),
		false,
	)
	if err != nil {
		log.Fatalf("[AGENT] Connect to Minio failed: %v\n", err)
	}

	// Create buckets
	buckets := []string{
		config.GetString(config.KEY_MINIO_VIDEO_BUCKET),
		config.GetString(config.KEY_MINIO_COVER_BUCKET),
		config.GetString(config.KEY_MINIO_AVATAR_BUCKET),
		config.GetString(config.KEY_MINIO_BACKGRAUND_BUCKET),
	}
	location := config.GetString("minio.location")
	retryTime := 0
	for {
		err := createBuckets(client, buckets, location)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				retryTime++
				log.Printf("[AGENT] create buckets timeout, retry %d times\n", retryTime)
			} else {
				log.Fatalf("[AGENT] create buckets fail: %v", err)
			}
			if retryTime > 10 {
				log.Fatalf("[AGENT] create buckets fail, retry %d times\n, the last error: %v", retryTime, err)
			}
		} else {
			break
		}
		time.Sleep(3 * time.Second)
	}

	// Connect to RabbitMQ
	var mq *rabbitmq.RabbitMQ
	retryTime = 0
	for {
		mq, err = rabbitmq.NewConsumer(
			config.GetString(config.KEY_RABBITMQ_HOST),
			config.GetString(config.KEY_RABBITMQ_PORT),
			config.GetString(config.KEY_RABBITMQ_USER),
			config.GetString(config.KEY_RABBITMQ_PASSWORD),
			config.GetString(config.KEY_RABBITMQ_QUEUE),
		)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				retryTime++
				log.Printf("[AGENT] connect to rabbitmq timeout, retry %d times\n", retryTime)
			} else {
				log.Fatalf("[AGENT] connect to rabbitmq fail: %v", err)
			}
			if retryTime > 10 {
				log.Fatalf("[AGENT] connect to rabbitmq fail, retry %d times\n, the last error: %v", retryTime, err)
			}
		} else {
			break
		}
		time.Sleep(3 * time.Second)
	}
	defer mq.Close()

	log.Println("[AGENT] agent is running")

	minioEndpoint := config.GetString(config.KEY_MINIO_ENDPOINT)
	coverBucket := config.GetString(config.KEY_MINIO_COVER_BUCKET)
	videoBucket := config.GetString(config.KEY_MINIO_VIDEO_BUCKET)
	JpegType := "image/jpeg"
	baseUrl := fmt.Sprintf("http://%s/%s/", minioEndpoint, videoBucket)
	for {
		filename := mq.Consume()
		url := baseUrl + filename + ".mp4"

		// get cover
		coverData, err := readFrameAsJpeg(url)
		if err != nil {
			log.Printf("[AGENT] read frame as jpeg failed: %v\n", err)
			continue
		}

		// upload cover
		coverReader := bytes.NewReader(coverData)
		_, err = client.PutObject(coverBucket, filename+".jpg", coverReader, coverReader.Size(), minio.PutObjectOptions{ContentType: JpegType})
		if err != nil {
			log.Printf("[AGENT] upload cover failed: %v\n", err)
			continue
		}
	}
}

func createBuckets(client *minio.Client, buckets []string, location string) error {
	for _, bucket := range buckets {
		exist, err := client.BucketExists(bucket)
		if err != nil {
			return err
		}
		if !exist {
			err := client.MakeBucket(bucket, location)
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

func readFrameAsJpeg(filePath string) ([]byte, error) {
	reader := bytes.NewBuffer(nil)
	err := ffmpeg.Input(filePath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(reader, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, img, nil)

	return buf.Bytes(), err
}
