package main

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func filesToGenerate(req *plugin.CodeGeneratorRequest) ([]*descriptor.FileDescriptorProto, error) {
	genFiles := make([]*descriptor.FieldDescriptorProto, 0)

	findMatchAndAppend := func(name string) error {
		for _, f := range req.ProtoFile {
			if f.GetName() == name {
				genFiles = append(genFiles, f)
				return nil
			}
		}

		return fmt.Errorf("could not find file named: %s", name)
	}

	for _, name := range req.FileToGenerate {
		if err := findMatchAndAppend(name); err != nil {
			return nil, err
		}
	}

	return genFiles, nil
}
