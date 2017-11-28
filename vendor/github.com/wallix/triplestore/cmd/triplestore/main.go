package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	tstore "github.com/wallix/triplestore"
)

var (
	outFormatFlag, inFormatFlag string
	baseFlag                    string
	dotPredicateFlag            string
	filesFlag                   arrayFlags
	prefixesFlag                arrayFlags
	useRdfPrefixesFlag          bool
)

func init() {
	flag.StringVar(&outFormatFlag, "out", "ntriples", "output format (ntriples, bin)")
	flag.StringVar(&inFormatFlag, "in", "bin", "input format (ntriples, bin)")
	flag.Var(&filesFlag, "files", "input file paths")
	flag.BoolVar(&useRdfPrefixesFlag, "rdf-prefixes", false, "use default RDF prefixes (rdf, rdfs, xsd)")
	flag.Var(&prefixesFlag, "prefix", "RDF custom prefixes (format: \"prefix:http://my.uri\"")
	flag.StringVar(&baseFlag, "base", "", "RDF custom base prefix")
	flag.StringVar(&dotPredicateFlag, "predicate", "", "Predicate on which to build a dot graph file")
}

func main() {
	flag.Parse()
	if len(filesFlag) == 0 {
		log.Fatal("need at list an argument `-files INPUT_FILE`")
	}
	context, err := buildContext(useRdfPrefixesFlag, prefixesFlag, baseFlag)
	if err != nil {
		log.Fatal(err)
	}
	if err := convert(filesFlag, outFormatFlag, context); err != nil {
		log.Fatal(err)
	}
}

func buildContext(useRdfPrefixes bool, prefixes []string, base string) (*tstore.Context, error) {
	var context *tstore.Context
	if useRdfPrefixes {
		context = tstore.RDFContext
	} else {
		context = tstore.NewContext()
	}
	for _, prefix := range prefixes {
		splits := strings.SplitN(prefix, ":", 2)
		if splits[0] == "" || splits[1] == "" {
			return context, fmt.Errorf("invalid prefix format: '%s'. expected \"prefix:http://my.uri\"", prefix)
		}
		context.Prefixes[splits[0]] = splits[1]
	}
	context.Base = base
	return context, nil
}

func convert(inFilePaths []string, outFormatFlag string, context *tstore.Context) error {
	var inFiles []io.Reader
	for _, inFilePath := range inFilePaths {
		in, err := os.Open(inFilePath)
		if err != nil {
			return fmt.Errorf("open input file '%s': %s", inFilePath, err)
		}
		inFiles = append(inFiles, in)
	}

	var inDecoder func(io.Reader) tstore.Decoder
	switch inFormatFlag {
	case "bin":
		inDecoder = tstore.NewBinaryDecoder
	case "ntriples":
		inDecoder = tstore.NewLenientNTDecoder
	default:
		return fmt.Errorf("unknown in flag '%s': expect 'ntriples' or 'bin'", outFormatFlag)
	}

	triples, err := tstore.NewDatasetDecoder(inDecoder, inFiles...).Decode()
	if err != nil {
		return err
	}

	var encoder tstore.Encoder
	switch outFormatFlag {
	case "ntriples":
		encoder = tstore.NewLenientNTEncoderWithContext(os.Stdout, context)
	case "bin":
		encoder = tstore.NewBinaryEncoder(os.Stdout)
	case "dot":
		if dotPredicateFlag == "" {
			return fmt.Errorf("missing -predicate param to output to dot format")
		}
		encoder = tstore.NewDotGraphEncoder(os.Stdout, dotPredicateFlag)
	default:
		return fmt.Errorf("unknown out flag '%s': expect 'ntriples, 'dot' or 'bin'", outFormatFlag)
	}

	if err := encoder.Encode(triples...); err != nil {
		return err
	}

	return nil
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
