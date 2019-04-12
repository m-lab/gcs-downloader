package main

import (
	"context"
	"flag"
	"log"
	"path"
	"regexp"
	"time"

	"cloud.google.com/go/storage"

	"github.com/m-lab/gcs-downloader/atomicfile"
	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/storagex"

	"github.com/m-lab/go/memoryless"
	"github.com/m-lab/go/prometheusx"
	"github.com/m-lab/go/rtx"
)

const sourceBucketPattern = `^gs://(?P<bucket>[^/]*)/(?P<path>.*)`

var (
	source string
	output string
	once   bool
	period time.Duration
)

func init() {
	log.SetFlags(log.LUTC | log.Lshortfile | log.LstdFlags)
	flag.StringVar(&source, "source", "", "gs://<bucket>[/<path>]")
	flag.DurationVar(&period, "period", time.Minute, "Schedule downloads to occur ever period on average.")
	flag.StringVar(&output, "output", ".", "Output directory name.")
	flag.BoolVar(&once, "once", false, "Download the source only once and then exit.")
}

var tempNew = atomicfile.New

func run(ctx context.Context, bucket storagex.Bucket, prefix string) error {
	return memoryless.Run(ctx, func() {
		ctx, cancel := context.WithTimeout(ctx, 2*period)
		defer cancel()
		bucket.Walk(ctx, prefix, func(o *storagex.Object) error {
			err := atomicfile.SaveFile(ctx, o, tempNew(path.Join(output, o.LocalName())))
			if err != nil {
				log.Println(err)
			}
			return nil
		})
	}, memoryless.Config{Expected: period, Once: once})
}

func main() {
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Failed to read args from env")

	srv := prometheusx.MustServeMetrics()
	defer srv.Close()
	fields := regexp.MustCompile(sourceBucketPattern).FindStringSubmatch(source)
	log.Printf("Copying %q to %q", source, output)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	rtx.Must(err, "Failed to create client")

	bucket := storagex.Bucket{BucketHandle: client.Bucket(fields[1])}
	err = run(ctx, bucket, fields[2])
	rtx.Must(err, "Failed to run the memoryless loop")
}
