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
)

func main() {
	flag.Parse()

	s := session.New(&aws.Config{Region: aws.String(*region)})
	kc := kinesis.New(s)

	streamName := aws.String(*stream)

	streams, err := kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot describe stream. please verify your stream is accecible.: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Your Stream information: %v\n", streams)

	iteratorOutput, err := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
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

	if err := iter(kc, iteratorOutput.ShardIterator, *maxItemSize); err != nil {
		fmt.Fprintf(os.Stderr, "get records failed: %s", err)
		os.Exit(1)
	}
}

// iter fetch records.
func iter(kc *kinesis.Kinesis, shardIterator *string, maxItemSize int) error {
	records, err := kc.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: shardIterator,
	})
	if err != nil {
		return err
	}
	for _, r := range records.Records {
		fmt.Printf("ApproximateArrivalTimestamp: %v\n", r.ApproximateArrivalTimestamp)
		if len(r.Data) > maxItemSize {
			fmt.Printf("Data: %s\n", r.Data[:maxItemSize-1])
		} else {
			fmt.Printf("Data: %s\n", r.Data[:])
		}
		fmt.Printf("SequenceNumber: %s\n", *r.SequenceNumber)
	}
	time.Sleep(time.Second * 3)
	return iter(kc, records.NextShardIterator, maxItemSize)
}
