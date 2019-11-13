# gcs-downloader
[![Version](https://img.shields.io/github/tag/m-lab/gcs-downloader.svg)](https://github.com/m-lab/gcs-downloader/releases) [![Build Status](https://travis-ci.org/m-lab/gcs-downloader.svg?branch=master)](https://travis-ci.org/m-lab/gcs-downloader) [![Coverage Status](https://coveralls.io/repos/m-lab/gcs-downloader/badge.svg?branch=master)](https://coveralls.io/github/m-lab/gcs-downloader?branch=master) [![GoDoc](https://godoc.org/github.com/m-lab/gcs-downloader?status.svg)](https://godoc.org/github.com/m-lab/gcs-downloader) [![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/gcs-downloader)](https://goreportcard.com/report/github.com/m-lab/gcs-downloader) 

A simple daemon for regularly downloading configurations from GCS.


Example usage:

```
docker run -i -t measurementlab/gcs-downloader:v0.1 -once \
  -source gs://fake-bucket-name/sub/path/
  -output .

```

