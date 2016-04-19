# kinesis-tail

Continuously fetch records from Kinesis Stream and print it. Useful for debugging Kinesis.

## How to use

	$ go get github.com/suzuken/kinesis-tail
	$ kinesis-tail -stream your-kinesis-stream-name -region ap-northeast-1 -max-item-size=512
	Your Stream information: {
	  StreamDescription: {
		HasMoreShards: false,
		RetentionPeriodHours: 24,
		Shards: [{
			HashKeyRange: {
			  EndingHashKey: "340282366920938463463374607431768211455",
			  StartingHashKey: "0"
			},
			SequenceNumberRange: {
			  StartingSequenceNumber: "49556705040321186153938184833717233207750712554199973890"
			},
			ShardId: "shardId-000000000000"
		  }],
		StreamARN: "arn:aws:kinesis:ap-northeast-1:xxxxxxxxxxxx:stream/your-kiness-stream-name",
		StreamName: "your-kiness-stream-name",
		StreamStatus: "ACTIVE"
	  }
	}
	{
	  ShardIterator: "AAAAAAAAAAF6jxGy30eyuTD4GzFrqAwEUFdWMIiWtFqpzGIShgXfo40VHDGHzNAdN0uvK5j6llWkEJQtdKZC3kQTIOQ/Y7Zj/lyBrV4GTkqz/hT+m/8TCZHU61fbEqU7Wd9V+Wa9szmeBQfdN5BAZDybhSUctjJRT+mgfbxVPXOjLkIrBNpAvzPFLRRaNSJj/ZwxIcG3k3nfwY7yUrcNb42OR6Cu7o1P"
	}
	ApproximateArrivalTimestamp: 2015-12-01 02:14:23 +0000 UTC
	Data: {"id":"5310e2ef-25e3-4942-87df-2aed28320b99","hoge":"fuga"}
	SequenceNumber: 49556846757453585699518614739447433851727187225060835330
	ApproximateArrivalTimestamp: 2015-12-01 02:24:41 +0000 UTC
	Data: {"id":"75ba2e4a-a68d-4d50-89cf-8b954a115a08","hoge":"fuga"}
	SequenceNumber: 49556846757453585699518614739460732035742990683338702850
	ApproximateArrivalTimestamp: 2015-12-01 04:10:58 +0000 UTC
	Data: {"id":"0d669fb6-045e-4299-946c-5786037ab638","hoge":"fuga"}
	...

## Usage

```
Usage of kineis-tail:
  -f    tailing kinesis stream forever or not (like: tail -f) (default true)
  -interval duration
        seconds for waiting next GetRecords request. (default 3s)
  -iterator-type string
        iterator type. Choose from TRIM_HORIZON(default), AT_SEQUENCE_NUMBER, or LATEST. (default "TRIM_HORIZON")
  -limit int
        limit records length of each GetRecords request. (default 100)
  -max-item-size int
        max byte size per item for printing. (default 4096)
  -region string
        the AWS region where your Kinesis Stream is. (default "ap-northeast-1")
  -stream string
        your stream name (default "your-stream")
```

## How does it works:

* At first, describe stream.
* Get shard iterator.
* `GetRecords` iteratively.

## LICENSE

MIT

## Author

Kenta Suzuki (a.k.a suzuken)
