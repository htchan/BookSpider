package main

//	@title			Book Spider Backend
//	@version		1.0
//	@description	This is a service fetching Novel from web

//go:generate go tool swag fmt -d .,../../internal
//go:generate go tool swag init -d ../../ -g cmd/api/swagger.go -o ../../docs -ot go,json
