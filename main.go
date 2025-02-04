package main

import (
	"movie-festival/cmd"
	"movie-festival/helpers"
)

func main() {

	// load config
	helpers.SetupConfig()

	// load log
	helpers.SetupLogger()

	// load db
	helpers.SetupMySQL()

	// run grpc
	go cmd.ServeGRPC()

	// run http
	cmd.ServeHTTP()
}
