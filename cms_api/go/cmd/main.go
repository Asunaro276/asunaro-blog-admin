package main

import (
	route "cms_api/internal/di"
)

func main() {
	e := route.RouteHandler()
	e.Logger.Fatal(e.Start(":8080"))
}
