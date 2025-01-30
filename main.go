package main

func main() {
	parseFlags()
	checkForFlagIncompatabilities()
	registerHandlers()
	startServer()
}
