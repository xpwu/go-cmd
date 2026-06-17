package cmd

func SetExeNameInUsageForTesting(name string) {
	exeNameInUsage = name
}

func SetOsExitForTesting(exit func(code int)) {
	osExit = exit
}

func ExitKeepaliveForTesting() {
	exitKeepAlive <- struct{}{}
}
