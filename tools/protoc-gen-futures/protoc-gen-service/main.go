package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/monzo/terrors"
)

const Version = "1.0.0"

func main() {
	versionFlag := flag.Bool("version", false, "print version & exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	_, err := readGenRequest(os.Stdin)
	if err != nil {
		errString := err.Error()
		writeResponse(os.Stdout, &plugin.CodeGeneratorResponse{Error: &errString})
		return
	}

}

func readGenRequest(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to read input files", nil)
	}

	req := new(plugin.CodeGeneratorRequest)
	if err := proto.Unmarshal(data, req); err != nil {
		return nil, terrors.Augment(err, "Failed to unmarshal input proto files", nil)
	}

	if len(req.FileToGenerate) == 0 {
		return nil, terrors.Augment(err, "No input proto files", nil)
	}

	return req, nil
}

func writeResponse(w io.Writer, rsp *plugin.CodeGeneratorResponse) {
	data, err := proto.Marshal(rsp)
	if err != nil {
		// If we can't reliably write an error in the proper way to the response; then lets give up and panic.
		log.Panicf("Failed to marshal response: %v", err)
	}

	if _, err := w.Write(data); err != nil {
		log.Panicf("Failed to write response: %v", err)
	}
}

func generateOutputFiles(in *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	return nil, nil
}
