## Fuzzing decoders

For context have a look at [go-fuzz](https://github.com/dvyukov/go-fuzz)

Corpus sample data is in directory `fuzz/corpus/samples.aws`

To fuzz the template decoding for instance do the following steps:

1. Build with

```sh
go-fuzz-build github.com/wallix/awless/template/fuzz/{parsing,parameters}
```

2. Then

```sh
go-fuzz -bin={parsing,parameters}-fuzz.zip -workdir=workdir
```

3. Stop (with Ctr+C) when enough. Look at the results. Fix the bugs and clean up the generated unneeded data (`rm -rf {parsing,parameters}/workdir/{crashers,suppressions}`)
