## ded

Find dead links in your source files.

### Contents

* [Install](#install)
* [Usage](#usage)
* [Examples](#examples)
* [Contribute](#contribute)
* [Credits](#credits)
* [License](#license)

### Install

#### Binaries

[View latest release.](https://github.com/kshvmdn/ded/releases/latest)

#### via Go

```sh
$ go get -u -v github.com/kshvmdn/ded/cmd/dep
$ which ded
$GOPATH/bin/ded
```

#### Build manually

```sh
$ git clone https://github.com/kshvmdn/ded.git $GOPATH/src/github.com/kshvmdn/ded
$ cd $_ # $GOPATH/src/github.com/kshvmdn/ded
$ go build -o ded -v github.com/kshvmdn/ded/cmd/ded
```

### Usage

* ded expects one or more files as input. It'll search these files for a list of
  valid URLs and report whether each URL resolves or not. View examples
  [below](#examples).

  ```sh
  $ ded file ...
  ...
  ```

* Use the `-help` flag to view the help dialogue.

  ```sh
  $ ded -help
  Usage of ded:
    -f string
          Output format (default "{{ .Filename }}:{{ .Line }}:{{ .Column }}: {{ .URL }} -> {{ .Result }}")
    -i string
          Regular expression to ignore
    -n    Hide "200 OK" responses
    -o string
          Output file or [stderr|stdout] (default "stdout")
    -v    Show version and exit
    -w int
          Number of workers (0 = number of CPUs)
  ```

### Examples

* Basic examples:

  ```sh
  $ ded main.go README.md
  ...
  $ ded ./**/*.go
  ...
  $ ls | xargs ded
  ...
  $ cat main.go | ded
  ...
  ```

* Custom output format:

  ```sh
  $ ded -f "{{ .URL }} {{ .Result }}" main.go README.md
  https://github.com/kshvmdn 200 OK
  http://nba.com/kashav-madan 404 Not Found
  ...
  ```

* Ignore a specific domain:

  ```sh
  $ ded -i "nba.com" main.go README.md
  main.go:95:4: https://github.com/kshvmdn -> 200 OK
  ...
  ```

* Hide any URL that 200s:

  ```sh
  $ ded -n main.go README.md
  README.md:47:16: http://nba.com/kashav-madan -> 404 Not Found
  ...
  ```

### Contribute

This project is completely open source, feel free to
[open an issue](https://github.com/kshvmdn/ded/issues) or
[submit a pull request](https://github.com/kshvmdn/ded/pulls).

### Credits

ded was inspired by [client9/misspell](https://github.com/client9/misspell).

### License

ded source code is released under the [MIT license](LICENSE).
