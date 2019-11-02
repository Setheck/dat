package cmd

var truePtr = true
var falsePtr = false

func ExampleVersion() {
	rc := RootCommand{
		ver: &truePtr,
	}
	rc.Init()
	rc.cmd.Run(nil, []string{"-v"})
	// Output: dat - version:v0.0.0 build:2019-11-02T01:23:46-0700
}

func ExampleWithEpoch() {
	rc := RootCommand{
		ver:   &falsePtr,
		all:   &falsePtr,
		local: &falsePtr,
		utc:   &falsePtr,
	}
	rc.Init()
	rc.cmd.Run(nil, []string{"1572683546"})
	// Output:
	// epoch: 1572683546
	// local: 11/02/2019 01:32:26 -0700
	//   utc: 11/02/2019 08:32:26 +0000
}
