# High Concurrency Geocoder

Creates a single output CSV using State of Indiana Public Geocoder using Goroutines for concurrency

**NOTE**: This provides no support for addresses outside of Indiana

# Usage

The implemetation requires an input csv with the following format:

| OBJECTID | Street | City | State | ZIP |
| --- | --- | --- | --- | --- |
| 1 | 1234 Main St | Indianapolis | IN | 46219 |

Set your output path in the ```config.yml``` file and (OPTIONAL) adjust the number of ```Goroutines``` to use.  Best performance has been profiled with the number of ```Goroutines``` matching the number of cores.

Run the geocoder

```shell
$ ./geocoder <path to csv>

```

# Notes

The compiled solution only requires that ```config.yml``` be present in the run directory.

# Building

Ensure conformity to ```Gofmt```

```shell
$ go install
$ go build

```

Build patterns

Linux:

```shell
$ env GOOS=linux GOARCH=amd64 go build geocoder.io/geocoder

```

Windows:

```shell
$ env GOOS=windows GOARCH=amd64 go build geocoder.io/geocoder

```

Mac

```shell
$ env GOOS=darwin GOARCH=amd64 go build geocoder.io/geocoder

```



