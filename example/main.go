package main

import (
	"../ftp"
	"fmt"
	"os"
	"io/ioutil"
	"bufio"
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
	// r, err := c.Retr("syslog.log")
	r, err := c.Retr("'omvssp.setup.jcl(setup)'")
	checkError(err)
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	str := string(data)
	fmt.Println("syslog.log: ", str)


	fmt.Println("[End]")
}

func SubmitJob() {
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

	// Issue command
	err = c.Cmd("TYPE A")
	checkError(err)
	fmt.Println("Issue command ASCII")

	//defer c.Logout()
	fmt.Println("Issue command SITE FILETYPE=JES JESLRECL=80")
	err = c.Cmd("SITE FILETYPE=JES JESLRECL=80")
	checkError(err)

	fmt.Println("Submit job")
	f, err := os.Open("./setup.jcl")
	checkError(err) 
	/*fmt.Println("Retr omvssp.setup.jcl(setup)")
	r, err := c.Retr("'omvssp.setup.jcl(setup)'")
	checkError(err)*/

	r := bufio.NewReader(f)
	// fmt.Println("Stor omvssp.setup.jcl(setup)")
	// err = c.Stor("'omvssp.setup.jcl(setuphy)'", r)
	// jobid, err := c.SubmitJob("'omvssp.setup.jcl(setuphy)'", r)
	jobid, err := c.SubmitJob(r)
	checkError(err)
	fmt.Println("JOBID: ", jobid)
	fmt.Println("[Done]")
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
	GetUnixFile()
	// GetPDS()
	SubmitJob()
}


