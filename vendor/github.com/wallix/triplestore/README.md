[![Build Status](https://api.travis-ci.org/wallix/triplestore.svg?branch=master)](https://travis-ci.org/wallix/triplestore)
[![Go Report Card](https://goreportcard.com/badge/github.com/wallix/triplestore)](https://goreportcard.com/report/github.com/wallix/triplestore)

# Triple Store

Triple Store is a library to manipulate RDF triples in a fast and fluent fashion.

RDF triples allow to represent any data and its relations to other data. It is a very versatile concept and is used in [Linked Data](https://en.wikipedia.org/wiki/Linked_data), graph databases and representation, simple directed graph, etc....

Here the RDF triples implementation follows along the [W3C RDF concepts](https://www.w3.org/TR/rdf11-concepts/). (**Note that blank nodes and reification are not implemented**.). More digestible info on [RDF Wikipedia](https://en.wikipedia.org/wiki/Resource_Description_Framework)

## Roadmap

- Generic map and filtering on triples
- High level triples query API
- RDF graph comparison
- Simple RDF graph traversal API
- Codec to [Turtle syntax](https://en.wikipedia.org/wiki/Turtle_(syntax))

## Triples quickstart

RDF is a resource description framework that allows to describe anything using triples. This is a powerful concept!

A triple simply consists of:

```
subject -> predicate -> object
```

... or you can also view that as: 

```
entity -> attribute -> value
```

## Library 

This library is written using the [Golang](https://golang.org) language.

Get it:

```sh
go get -u github.com/wallix/triplestore
```

Test it:

```
go test -v -cover -race github.com/wallix/triplestore
```

Import it in your source code:

```go
import (
	"github.com/wallix/triplestore"
	// tstore "github.com/wallix/triplestore" for less verbosity
)
```

## Usage

### Manipulating triples

```go
triples = append(triples,
	SubjPred("me", "name").StringLiteral("jsmith"),
 	SubjPred("me", "age").IntegerLiteral(26),
 	SubjPred("me", "male").BooleanLiteral(true),
 	SubjPred("me", "born").DateTimeLiteral(time.Now()),
 	SubjPRed("me", "mother").Resource("mum#121287"),
)
```

Although you can build triples the way you want to model any data, they are usually built from known RDF vocabularies & namespace. Ex: [foaf](http://xmlns.com/foaf/spec/), ...

Check if triples are equal:

```go
	me := SubjPred("me", "name").StringLiteral("jsmith")
 	you := SubjPred("me", "name").StringLiteral("fdupond")

 	if me.Equal(you) {
 	 	...
 	}
)
```

### Triple Storage

// TODO

### Triple Queries

// TODO

### RDF Graph

// TODO

### Raw storage

In this library, higher level APIs used `Encoders` and `Decoders` to store and exchanges triples. Triples can be encoded & decoded using a simple format in order to persists, flush or send them over the network.

For example

```go
triples = append(triples,
	SubjPred("me", "name").StringLiteral("jsmith"),
	...
 	SubjPred("me", "born").DateTimeLiteral(time.Now()),
)

enc := NewBinaryEncoder(myWriter)
err := enc.Encode(triples)
...

dec := NewBinaryDecoder(myReader)
triples, err := dec.Decode()

```
