package main

import (
	"time"
	"log"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type StreamItem struct {
	TableName string
	Data      map[string]interface{}
	ProcessTime time.Time
}

type Streamer interface {
	Stream(*StreamItem) error
}

type KinesisStreamer struct {
	Config *Config
}

func (s *KinesisStreamer) Stream(item *StreamItem) error {

	log.Printf("Streaming item %+v", item)
	return nil
}

func NewStreamer(c *Config) *KinesisStreamer {

	for _,table := range c.IncludeTables{
		_, err := c.Kinesis.DescribeStreamSummary(&kinesis.DescribeStreamSummaryInput{
			StreamName: &table,
		})

		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == kinesis.ErrCodeResourceNotFoundException {
			log.Fatalf("the following kinesis stream does not exists %s", table)
		}
	}

	return &KinesisStreamer{
		Config: c,
	}
}

