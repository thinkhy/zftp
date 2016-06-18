# zftp
Golang client for z/OS FTP server

zftp is forked from github.com/jlaffaye/ftp and makes some changes to accommdate to z/OS FTP
server

Currently zftp only supports z/OS FTP server with JESINTERFACELEVEL=2

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
	// TODO
}

```

Reference
----------------
   * [z/OS V2R1.0 Communications Server: IP User's Guide and Commands](http://publibz.boulder.ibm.com/epubs/pdf/f1a2b900.pdf)  	


