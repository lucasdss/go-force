# go-force

[![Go Report Card](https://goreportcard.com/badge/github.com/taxnexus/go-force)](https://goreportcard.com/report/github.com/taxnexus/go-force)
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/taxnexus/go-force/force)

[Golang](http://golang.org/) API wrapper for Salesforce REST and Streaming APIs

This is a fork of an older package principally written by earlier contributors that has been enhanced to support Push topics in the Salesforce Streaming API.

## Installation

    go get github.com/taxnexus/go-force/force

## Example

```golang
package main

import (
	"fmt"
	"log"

	"github.com/taxnexus/go-force/force"
	"github.com/taxnexus/go-force/sobjects"
)

type someCustomSObject struct {
	sobjects.BaseSObject

	Active    bool   `force:"Active__c"`
	AccountID string `force:"Account__c"`
}

func (t *someCustomSObject) apiName() string {
	return "SomeCustomObject__c"
}

type someCustomSObjectQueryResponse struct {
	sobjects.BaseQuery

	Records []*SomeCustomSObject `force:"records"`
}

func main() {
	// Init the force
	forceAPI, err := force.Create(
		"YOUR-API-VERSION",
		"YOUR-CLIENT-ID",
		"YOUR-CLIENT-SECRET",
		"YOUR-USERNAME",
		"YOUR-PASSWORD",
		"YOUR-SECURITY-TOKEN",
		"YOUR-ENVIRONMENT",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get somCustomSObject by ID
	someCustomSObject := &SomeCustomSObject{}
	err = forceAPI.GetSObject("Your-Object-ID", someCustomSObject)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", someCustomSObject)

	// Query
	someCustomSObjects := &SomeCustomSObjectQueryResponse{}
	err = forceAPI.Query("SELECT Id FROM SomeCustomSObject__c LIMIT 10", someCustomSObjects)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", someCustomSObjects)
}
```

## Documentation

- [Package Reference](http://godoc.org/github.com/taxnexus/go-force/force)
- [Force.com API Reference](http://www.salesforce.com/us/developer/docs/api_rest/)
