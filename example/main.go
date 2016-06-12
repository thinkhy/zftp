package main

import (
	"../ftp.go"
	"os"
)

func main() {
	server := os.Getenv("TEST_FTP_SERVER")
	ftp := new(ftp.FTP)
	ftp.Debug = true
	ftp.Connect(server, 21)

}
