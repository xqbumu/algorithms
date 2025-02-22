package main

import (
	"fmt"

	"github.com/apache/opendal-go-services/memory"
	"github.com/apache/opendal-go-services/s3"
	opendal "github.com/apache/opendal/bindings/go"
)

func main() {
	RunS3()
}

func RunS3() {
	// Initialize a new in-s3 operator
	op, err := opendal.NewOperator(s3.Scheme, opendal.OperatorOptions{
		"root":              "/opendal",
		"endpoint":          "http://s3.minio.lan",
		"bucket":            "playground",
		"region":            "lan-infra",
		"access_key_id":     "aPQ6s7raoqK70WkyADsz",
		"secret_access_key": "jmVFx1cdrFgwNIX6I0ATyYD1edxbMcjRzdjpVlNu",
	})
	if err != nil {
		panic(err)
	}
	defer op.Close()

	// Write data to a file named "test"
	err = op.Write("hello.txt", []byte("Hello opendal go binding!"))
	if err != nil {
		panic(err)
	}
}

func RunMemory() {
	// Initialize a new in-memory operator
	op, err := opendal.NewOperator(memory.Scheme, opendal.OperatorOptions{})
	if err != nil {
		panic(err)
	}
	defer op.Close()

	// Write data to a file named "test"
	err = op.Write("test", []byte("Hello opendal go binding!"))
	if err != nil {
		panic(err)
	}

	// Read data from the file "test"
	data, err := op.Read("test")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read content: %s\n", data)

	// List all entries under the root directory "/"
	lister, err := op.List("/")
	if err != nil {
		panic(err)
	}
	defer lister.Close()

	// Iterate through all entries
	for lister.Next() {
		entry := lister.Entry()

		// Get entry name (not used in this example)
		_ = entry.Name()

		// Get metadata for the current entry
		meta, _ := op.Stat(entry.Path())

		// Print file size
		fmt.Printf("Size: %d bytes\n", meta.ContentLength())

		// Print last modified time
		fmt.Printf("Last modified: %s\n", meta.LastModified())

		// Check if the entry is a directory or a file
		fmt.Printf("Is directory: %v, Is file: %v\n", meta.IsDir(), meta.IsFile())
	}

	// Check for any errors that occurred during iteration
	if err := lister.Error(); err != nil {
		panic(err)
	}

	// Copy a file
	op.Copy("test", "test_copy")

	// Rename a file
	op.Rename("test", "test_rename")

	// Delete a file
	op.Delete("test_rename")
}
