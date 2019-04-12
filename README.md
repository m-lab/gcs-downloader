# gcs-downloader

A simple daemon for regularly downloading configurations from GCS.


Example usage:

```
docker run -i -t measurementlab/gcs-downloader:v0.1 -once \
  -source gs://fake-bucket-name/sub/path/
  -output .

```
