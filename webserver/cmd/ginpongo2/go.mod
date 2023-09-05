module github.com/telemac/goutils/webserver/cmd/ginpongo2

go 1.16

require (
	github.com/flosch/pongo2 v0.0.0-20200913210552-0d938eb266f3
	github.com/gin-gonic/gin v1.8.1
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/tecome/pongo2gin v0.0.0-20210616004019-cf499344efb9
	github.com/telemac/goutils v0.0.0-00010101000000-000000000000
)

replace github.com/telemac/goutils => ../../../
