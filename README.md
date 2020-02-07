# search-cli

## Setup
To build an executable, do the following: 

1. Make sure you have at least go version 1.13 installed (https://golang.org/dl/)
2. Make sure go modules are enabled (if you have the latest version of go, they are enabled by default)
3. Clone this repository
4. Run the unit tests. This will automatically install all dependencies. 
5. Run `build.sh` to build an executible at `bin/zensearch`

## Testing
Run the following at the root of the project
```
$ go test ./...
```

## Command Line Docs
You can run the following command to view details on how to use the program:
```
$ ./bin/zensearch help
```

[You can also find documentation in this repository](docs/zensearch.md)
