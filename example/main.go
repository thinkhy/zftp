package main

import (
	"../ftp"
	"fmt"
	"os"
	"io/ioutil"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func GetUnixFile() {
	// thinkTime := 3000 * time.Millisecond
	server := os.Getenv("TEST_FTP_SERVER")
	user := os.Getenv("TEST_FTP_USER")
	password := os.Getenv("TEST_FTP_PASSWORD")
	addr := fmt.Sprintf("%s:21", server)

	// dail
	c, err := ftp.DialTimeout(addr, 5*time.Second)
	checkError(err)
	fmt.Println("FTP Connected")
	defer c.Quit()

	// login
	err = c.Login(user, password)
	checkError(err)
	//defer c.Logout()

        // Change dir
	err = c.ChangeDir("/tmp/")
	checkError(err)

	// PWD
	pwd, err := c.CurrentDir()
	checkError(err)
	fmt.Println("PWD: ", pwd)

	// Issue command
	err = c.Cmd("TYPE A")
	checkError(err)
	fmt.Println("Issue command ASCII")

        // Retrieve file
	r, err := c.Retr("syslog.log")
	checkError(err)
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	str := string(data)
	fmt.Println("syslog.log: ", str)

	fmt.Println("[End]")
}

func GetPDS() {
	// thinkTime := 3000 * time.Millisecond
	server := os.Getenv("TEST_FTP_SERVER")
	user := os.Getenv("TEST_FTP_USER")
	password := os.Getenv("TEST_FTP_PASSWORD")
	addr := fmt.Sprintf("%s:21", server)

	// dail
	c, err := ftp.DialTimeout(addr, 5*time.Second)
	checkError(err)
	fmt.Println("FTP Connected")
	defer c.Quit()

	// login
	err = c.Login(user, password)
	checkError(err)
	//defer c.Logout()


        // Change dir
	//err = c.ChangeDir("..")
	//err = c.ChangeDir("..")
	//err = c.ChangeDir("..")
	// err = c.ChangeDir("OMVSSP.SETUP.JCL")
	//checkError(err)
	//fmt.Println("Change DIR to OMVSSP.SETUP.JCL")

	err = c.ChangeDir("..")
	err = c.ChangeDir("..")
	err = c.ChangeDir("..")
	// err = c.ChangeDir("OMVSSP.SETUP.JCL")
	err = c.ChangeDir("/tmp")
	checkError(err)
	fmt.Println("Change DIR to OMVSSP.SETUP.JCL")
	entries, err := c.List("") 
	checkError(err)
	for _, et := range entries {
		fmt.Println("Name: ", et.Name)
	}
	
}

func main() {
	// GetUnixFile()
	GetPDS()
}
