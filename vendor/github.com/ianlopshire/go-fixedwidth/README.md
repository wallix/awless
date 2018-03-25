# fixedwidth [![GoDoc](https://godoc.org/github.com/ianlopshire/go-fixedwidth?status.svg)](http://godoc.org/github.com/ianlopshire/go-fixedwidth) [![Report card](https://goreportcard.com/badge/github.com/ianlopshire/go-fixedwidth)](https://goreportcard.com/report/github.com/ianlopshire/go-fixedwidth) [![Go Cover](http://gocover.io/_badge/github.com/ianlopshire/go-fixedwidth)](http://gocover.io/github.com/ianlopshire/go-fixedwidth)

Package fixedwidth provides encoding and decoding for fixed-width formatted Data.

`go get github.com/ianlopshire/go-fixedwidth`

## Usage

### Struct Tags
Position within a line is controlled via struct tags.
The tags should be formatted as `fixed:"{startPos},{endPos}"` where `startPos` and `endPos` are both positive integers greater than 0.
Positions start at 1. The interval is inclusive. Fields without tags are ignored.

### Encode
```go
// define some data to encode
people := []struct {
    ID        int     `fixed:"1,5"`
    FirstName string  `fixed:"6,15"`
    LastName  string  `fixed:"16,25"`
    Grade     float64 `fixed:"26,30"`
}{
    {1, "Ian", "Lopshire", 99.5},
}

data, err := fixedwidth.Marshal(people)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s", data)
// Output:
// 1    Ian       Lopshire  99.50
```

### Decode
```go
// define the format
var people []struct {
    ID        int     `fixed:"1,5"`
    FirstName string  `fixed:"6,15"`
    LastName  string  `fixed:"16,25"`
    Grade     float64 `fixed:"26,30"`
}

// define some fixed-with data to parse
data := []byte("" +
    "1    Ian       Lopshire  99.50" + "\n" +
    "2    John      Doe       89.50" + "\n" +
    "3    Jane      Doe       79.50" + "\n")


err := fixedwidth.Unmarshal(data, &people)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%+v\n", people[0])
fmt.Printf("%+v\n", people[1])
fmt.Printf("%+v\n", people[2])
// Output:
//{ID:1 FirstName:Ian LastName:Lopshire Grade:99.5}
//{ID:2 FirstName:John LastName:Doe Grade:89.5}
//{ID:3 FirstName:Jane LastName:Doe Grade:79.5}
```

## Licence
MIT
