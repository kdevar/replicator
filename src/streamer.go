package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"log"
	"time"
)

type StreamItem struct {
	TableName   string
	Data        map[string]interface{}
	ProcessTime time.Time
}

type Streamer interface {
	Stream(*StreamItem) error
}

type KinesisStreamer struct {
	Config *Config
}

func (s *KinesisStreamer) Stream(item *StreamItem) error {
	json, err := json.Marshal(item.Data)

	partitionKey := item.TableName

	if err != nil {
		log.Printf("%v", err)
		return err
	}

	out, err := s.Config.Kinesis.PutRecord(&kinesis.PutRecordInput{
		StreamName:   &item.TableName,
		Data:         json,
		PartitionKey: &partitionKey,
	})

	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("Streaming item=%+v output=%v", item, out)

	return nil
}

func NewStreamer(c *Config) *KinesisStreamer {

	for _, table := range c.IncludeTables {
		_, err := c.Kinesis.DescribeStreamSummary(&kinesis.DescribeStreamSummaryInput{
			StreamName: &table,
		})

		if isStreamNotFoundError(err) {
			log.Fatalf("the following kinesis stream does not exists %s", table)
		}
	}

	return &KinesisStreamer{
		Config: c,
	}
}

func isStreamNotFoundError(err error) bool {
	awsErr, ok := err.(awserr.Error)
	return ok && awsErr.Code() == kinesis.ErrCodeResourceNotFoundException
}
