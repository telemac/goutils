module github.com/telemac/goutils/tcpserver/sample

go 1.18

require github.com/telemac/goutils v1.1.28

require (
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/tevino/abool v1.2.0 // indirect
	golang.org/x/net v0.0.0-20220526153639-5463443f8c37 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)

replace (
	github.com/telemac/goutils  => ../../../goutils
)
