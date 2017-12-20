## Fuzzing decoders

For context have a look at [go-fuzz](https://github.com/dvyukov/go-fuzz)

Corpus sample data are in directories `fuzz/{ntriples,binary}/corpus/samples.*`

For instance, to fuzz the ntriples decoding for instance do the following steps:

1. Build with

```sh
go-fuzz-build github.com/wallix/triplestore/fuzz/ntriples
```

2. Then

```sh
go-fuzz -bin=ntriples-fuzz.zip -workdir=fuzz/ntriples/corpus
```

3. Stop (with Ctr+C) when enough. Look at the results. Fix the bugs and clean up the generated unneeded data (`rm -rf fuzz/ntriples/corpus/{corpus,crashers,suppressions}`)
