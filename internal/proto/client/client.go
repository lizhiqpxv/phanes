package client

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"phanes/internal/base"

	"github.com/spf13/cobra"
)

// CmdClient represents the source command.
var CmdClient = &cobra.Command{
	Use:   "client",
	Short: "Generate the proto client code",
	Long:  "Generate the proto client code. Example: kratos proto client helloworld.proto",
	Run:   run,
}

var protoPath string

func init() {
	if protoPath = os.Getenv("PROTO_PATH"); protoPath == "" {
		protoPath = "./third_party"
	}
	CmdClient.Flags().StringVarP(&protoPath, "proto_path", "p", protoPath, "proto path")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the proto file or directory")
		return
	}
	var (
		err   error
		proto = strings.TrimSpace(args[0])
	)
	if err = look("protoc-gen-go", "protoc-gen-go-grpc", "protoc-gen-micro"); err != nil {
		return
	}
	if strings.HasSuffix(proto, ".proto") {
		err = generate(proto, args)
	} else {
		err = walk(proto, args)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" || strings.HasPrefix(path, "third_party") {
			return nil
		}
		return generate(path, args)
	})
}

// generate is used to execute the generate command for the specified proto file
func generate(proto string, args []string) error {
	input := []string{
		"--proto_path=.",
	}
	if pathExists(protoPath) {
		input = append(input, "--proto_path="+protoPath)
	}
	inputExt := []string{
		"--proto_path=" + base.KratosMod(),
		"--proto_path=" + filepath.Join(base.KratosMod(), "third_party"),
		"--go_out=paths=source_relative:.",
		"--micro_out=paths=source_relative:.",
	}
	input = append(input, inputExt...)
	protoBytes, err := os.ReadFile(proto)
	if err == nil && len(protoBytes) > 0 {
		if ok, _ := regexp.Match(`\n[^/]*(import)\s+"validate/validate.proto"`, protoBytes); ok {
			input = append(input, "--validate_out=lang=go,paths=source_relative:.")
		}
	}
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
