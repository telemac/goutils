module github.com/telemac/goutils/tcpserver/sample

go 1.21

toolchain go1.22.0

require (
	github.com/sirupsen/logrus v1.9.3
	github.com/telemac/goutils v1.1.46
)

require (
	github.com/tevino/abool v1.2.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/telemac/goutils => ../../../goutils
