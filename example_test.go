package autocomplete_test

import (
	"flag"
	"os"

	"github.com/serge-v/autocomplete"
)

func ExampleHandleArgs() {
	flag.Parse()
	autocomplete.HandleArgs()
}

func ExampleHandle() {
	getFruits := func() []string {
		return []string{"apple", "pear", "kiwi"}
	}

	flag.Parse()
	autocomplete.Handle("fruit", getFruits)
	autocomplete.HandleArgs()
}
