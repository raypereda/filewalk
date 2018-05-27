package main_test

import (
	"github.com/raypereda/filewalk"
)

func Example() {
	main.Main()
	// Output:
	// 	51

	// 	# extension
	//    33
	//    10 .sample
	// 	2 .go
	// 	2 .dll
	// 	1 .gitignore
	// 	1 .exe
	// 	1 .idx
	// 	1 .pack
}
