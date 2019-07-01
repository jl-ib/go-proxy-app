package main
import (
	handlers "github.com/jl-ib/proxy-app/api/handlers"
	utils "github.com/jl-ib/proxy-app/api/utils"
	server "github.com/jl-ib/proxy-app/api/server"
)

func main() {
	/*
	Router Iris
	Env Vars
	*/

	utils.LoadEnv()
	app := server.SetUp()
	// middleware.InitQueue()
	handlers.HandlerRedirection(app)
	server.RunServer(app)
}