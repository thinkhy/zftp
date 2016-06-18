package zftp

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	//"bytes"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func TestPutGetDeleteUnixFile(t *testing.T) {
	server := os.Getenv("TEST_FTP_SERVER")
	user := os.Getenv("TEST_FTP_USER")
	password := os.Getenv("TEST_FTP_PASSWORD")
	// psdataset := os.Getenv("TEST_FTP_DATASET_PS")
	addr := fmt.Sprintf("%s:21", server)
	fmt.Println(addr)
	z, err := Dial(addr, 30)
	assert.Nil(t, err)
	err = z.Login(user, password)
	assert.Nil(t, err)

	local := "zftp_test.go"
	remote := "/tmp/testfile"
	f, err := os.Open(local)
	assert.Nil(t, err)
	r := bufio.NewReader(f)
	err = z.PutUnixFile(r, remote)
	assert.Nil(t, err)
	f.Close()
	rc, err := z.GetUnixFile(remote)
	assert.Nil(t, err)
	data, err := ioutil.ReadAll(rc)
	assert.Nil(t, err)
	rc.Close()
	str := string(data)
	data1, err := ioutil.ReadFile(local)
	assert.Equal(t, string(data1), str)
	err = z.DeleteUnixFile(remote)
	assert.Nil(t, err)
	_, err = z.GetUnixFile(remote)
	assert.NotNil(t, err)
	err = z.Quit()
	assert.Nil(t, err)
}

func TestGetPsDataset(t *testing.T) {
	server := os.Getenv("TEST_FTP_SERVER")
	user := os.Getenv("TEST_FTP_USER")
	password := os.Getenv("TEST_FTP_PASSWORD")
	psdataset := os.Getenv("TEST_FTP_DATASET_PS")
	addr := fmt.Sprintf("%s:21", server)
	fmt.Println(addr)
	z, err := Dial(addr, 30)
	assert.Nil(t, err)
	err = z.Login(user, password)
	assert.Nil(t, err)
	err = z.SetSeqMode()
	assert.Nil(t, err)
	r, err := z.GetPsDataset(psdataset)
	assert.Nil(t, err)
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.True(t, len(string(data)) > 0)
	err = z.Quit()
	assert.Nil(t, err)
}

func TestGetPdsDataset(t *testing.T) {
}

func TestSubmitJob(t *testing.T) {
	server := os.Getenv("TEST_FTP_SERVER")
	user := os.Getenv("TEST_FTP_USER")
	password := os.Getenv("TEST_FTP_PASSWORD")
	addr := fmt.Sprintf("%s:21", server)
	fmt.Println(addr)
	z, err := Dial(addr, 30)
	assert.Nil(t, err)
	err = z.Login(user, password)
	assert.Nil(t, err)
	err = z.SetJesMode()
	assert.Nil(t, err)
	// Do not use an end-of-line sequence other than CRLF if the server is a z/OS FTP server.
	// The z/OS FTP server supports only the CRLF("\r\n") value for incoming ASCII data.
	jcl :=
		`
//HELLOWLD JOB 'TT',MSGLEVEL=(1,1),MSGCLASS=H,CLASS=A,USER=MEGA 
//HELLOWLD EXEC PGM=IKJEFT01 
//SYSTSPRT DD SYSOUT=* 
//SYSTSIN DD * 
BPXBATCH SH echo "hello world"
/*
`
	jcl = z.Unix2Dos(jcl)
	r := strings.NewReader(jcl)
	fmt.Println("Submit job")
	jobid, err := z.SubmitJob(r)
	assert.Nil(t, err)
	fmt.Println("jobid:", jobid)
	fmt.Println("Purge job")
	j, err := z.GetJobStatusByID(jobid)
	assert.Nil(t, err)
	fmt.Println("job entry:", j)

	rc, err := z.GetJobLog(jobid)
	assert.Nil(t, err)
	data, err := ioutil.ReadAll(rc)
	assert.Nil(t, err)
	rc.Close()
	assert.True(t, strings.Contains(string(data), `BPXBATCH SH echo "hello world"`))

	err = z.PurgeJob(jobid)
	assert.Nil(t, err)
	err = z.Quit()
	assert.Nil(t, err)
}

func TestSubmitRemoteJob(t *testing.T) {
	server := os.Getenv("TEST_FTP_SERVER")
	user := os.Getenv("TEST_FTP_USER")
	password := os.Getenv("TEST_FTP_PASSWORD")
	psdataset := os.Getenv("TEST_FTP_DATASET_PS")
	addr := fmt.Sprintf("%s:21", server)
	fmt.Println(addr)
	z, err := Dial(addr, 30)
	assert.Nil(t, err)
	err = z.Login(user, password)
	assert.Nil(t, err)
	// err = z.SetJesMode()
	// assert.Nil(t, err)
	fmt.Println("Submit job")
	jobid, err := z.SubmitRemoteJob(psdataset)
	assert.Nil(t, err)
	fmt.Println("jobid:", jobid)
	fmt.Println("Purge job")
	j, err := z.GetJobStatusByID(jobid)
	assert.Nil(t, err)
	fmt.Println("job entry:", j)

	rc, err := z.GetJobLog(jobid)
	assert.Nil(t, err)
	data, err := ioutil.ReadAll(rc)
	assert.Nil(t, err)
	rc.Close()
	assert.True(t, strings.Contains(string(data), `BPXBATCH SH echo "hello world"`))

	err = z.PurgeJob(jobid)
	assert.Nil(t, err)
	err = z.Quit()
	assert.Nil(t, err)
}
