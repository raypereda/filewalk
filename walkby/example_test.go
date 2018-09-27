package walkby_test

import (
	"os"

	"github.com/raypereda/filewalk/walkby"
)

func ExampleWalkByExt() {
	os.Args = []string{"filewalk", "."}
	walkby.Main()
	// Output:
	//    # extension
	//    2 .go
}

func ExampleWalkByApp() {
	os.Args = []string{"filewalk", "-app", "."}
	walkby.Main()
	// Fix(raypereda): the output looks identical
	// Output:
	// banned count: 0
	// file count: 3
	//
}
