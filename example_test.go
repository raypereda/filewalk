package main_test

func ExampleWalkByExt() {
	H
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

func ExampleWalkByProject() {
	// TODO: fix this
	// no additional command-line parameters for this
	// main.Main()
	//  Output:
	// 	app: hss/App1
	// 	# extension
	// 	1 .pdf
	//  app: hss/App2
	// 	# extension
	// 	1 .sql
	//  app: hss/COE/app3
	// 	# extension
	// 	2 .pdf
	// 	1 .sql
	//  file count: 5
}
