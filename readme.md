# Robloxgo

[![Go Reference](https://pkg.go.dev/badge/github.com/RhykerWells/robloxgo.svg)](https://pkg.go.dev/github.com/RhykerWells/robloxgo) [![Go Report Card](https://goreportcard.com/badge/github.com/RhykerWells/robloxgo)](https://goreportcard.com/report/github.com/RhykerWells/robloxgo)

RobloxGo is a [Go](https://golang.org/) package that provides low level 
bindings to the [ROBLOX](https://create.roblox.com/docs/cloud) API. RobloxGo's current aims are support for specific functions relating to users & groups, but if the package gains enough traction other aspects will be included

**For help with this package, please use the [discussions](https://github.com/RhykerWells/robloxgo/discussions) tab.**

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` *will always pull the latest tagged release from the master branch.*

```sh
go get github.com/RhykerWells/robloxgo
```

### Usage

Import the package into your project.

```go
import "github.com/RhykerWells/robloxgo"
```

Construct a new Roblox client which is used to access the variety of 
Roblox API functions 

```go
client, err := robloxgo.Create("Your roblox API key")
```

## Documentation

**NOTICE**: This library and the ROBLOX API are unfinished.
Because of that there may be major changes to library in the future.

The Robloxgo code is currently the only documentation available.
Go reference (below) presents that information in a nice format.

- [![Go Reference](https://pkg.go.dev/badge/github.com/RhykerWells/robloxgo.svg)](https://pkg.go.dev/github.com/RhykerWells/robloxgo) 

## Contributing
Contributions are very welcomed, however please follow the below guidelines.

- Try to match current naming conventions as closely as possible.  
- This package is intended to be a low level direct mapping of the Roblox API, 
so please avoid adding enhancements outside of that scope without first 
discussing it.
- Create a Pull Request with your changes against the master branch.