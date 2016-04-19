package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"os"
	"time"
)

var (
	stream       = flag.String("stream", "your-stream", "your stream name")
	region       = flag.String("region", "ap-northeast-1", "the AWS region where your Kinesis Stream is.")
	iteratorType = flag.String("iterator-type", "TRIM_HORIZON", "iterator type. Choose from TRIM_HORIZON(default), AT_SEQUENCE_NUMBER, or LATEST.")
	maxItemSize  = flag.Int("max-item-size", 4096, "max byte size per item for printing.")
	forever      = flag.Bool("f", true, "tailing kinesis stream forever or not (like: tail -f)")
	limit        = flag.Int64("limit", 100, "limit records length of each GetRecords request.")
	interval     = flag.Duration("interval", time.Second*3, "seconds for waiting next GetRecords request.")
)

type Client struct {
	*kinesis.Kinesis
	maxItemSize int
	Limit       *int64
}

func main() {
	flag.Parse()

	s := session.New(&aws.Config{Region: aws.String(*region)})
	c := &Client{Kinesis: kinesis.New(s), maxItemSize: *maxItemSize, Limit: limit}

	streamName := aws.String(*stream)

	streams, err := c.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot describe stream. please verify your stream is accecible.: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Your Stream information: %v\n", streams)

	iteratorOutput, err := c.GetShardIterator(&kinesis.GetShardIteratorInput{
		// take first shard.
		ShardId:           streams.StreamDescription.Shards[0].ShardId,
		ShardIteratorType: aws.String(*iteratorType),
		StreamName:        streamName,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get iterator: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", iteratorOutput)

	iter := iteratorOutput.ShardIterator
	for {
		if iter == nil && *forever == false {
			fmt.Fprintf(os.Stderr, "all data in this stream is fetched. : %s")
			os.Exit(0)
		}
		iter, err = c.fetch(iter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get records failed: %s", err)
			os.Exit(1)
		}
		time.Sleep(*interval)
	}
}

// iter fetch records.
func (c *Client) fetch(shardIterator *string) (nextShardIterator *string, err error) {
	records, err := c.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: shardIterator,
		Limit:         c.Limit,
	})
	if err != nil {
		return nil, err
	}
	for _, r := range records.Records {
		fmt.Printf("ApproximateArrivalTimestamp: %v\n", r.ApproximateArrivalTimestamp)
		if len(r.Data) > c.maxItemSize {
			fmt.Printf("Data: %s\n", r.Data[:c.maxItemSize-1])
		} else {
			fmt.Printf("Data: %s\n", r.Data[:])
		}
		fmt.Printf("SequenceNumber: %s\n", *r.SequenceNumber)
	}
	return records.NextShardIterator, nil
}
