# dead-link-detector

Find dead links in your source files.

## Usage

  ```sh
  $ git clone https://github.com/kashav/go-dead-link-detector
  $ go build ./cmd/detector/main.go
  $ ./detector -help
  Usage of ./detector:
    -f string
        Output format (default "{{ .Filename }}:{{ .Line }}:{{ .Column }}: {{ .URL }} -> {{ .Result }}")
    -i string
        Regular expression to ignore
    -n	Hide "200 OK" responses
    -o string
        Output file or [stderr|stdout] (default "stdout")
    -v	Show version and exit
    -w int
        Number of workers (0 = number of CPUs)
  ```
