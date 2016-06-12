# zftp
Golang client for z/OS FTP server

zftp is forked from github.com/jlaffaye/ftp and makes some changes to accommdate to z/OS FTP
server

Install
-----------

```shell
go get github.com/thinkhy/zftp
```

Example
-----------

```golang
package main

import (
        "fmt"
        "os"
        "github.com/thinkhy/zftp"
       )

func main() {
   ftp := new(ftp.FTP)
   // debug default false
   ftp.Debug = true
   ftp.Connect("localhost", 21)

   // login
   ftp.Loign("user", "password")
   if ftp.Code == 530 {
       fmt.Println("error: login failure")
       os.Exit(-1)
   }

   // pwd
   ftp.Pwd()
   fmt.Println("code:", ftp.Code, ", message:", ftp.Message)

   // make dir
   ftp.Mkd("/path")
   ftp.Request("TYPE I")

   // stor file
   b, _ := ioutil.ReadFile("/path/a.txt")
   ftp.Stor("/path/a.txt", b)

   // TODO - 2016-06-12

   ftp.Quit()
}

Reference
----------------
   * [z/OS V2R1.0 Communications Server: IP User's Guide and Commands](http://publibz.boulder.ibm.com/epubs/pdf/f1a2b900.pdf)  	


```
