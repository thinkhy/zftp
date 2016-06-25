package zftp

import (
	"bufio"
	"bytes"
	ftp "github.com/thinkhy/zftp/ftp"
	"io"
	"io/ioutil"
	"time"
	// "strings"
	"fmt"
	"regexp"
)

type Zftp struct {
	*ftp.ServerConn
}

/* JOBNAME  JOBID    OWNER    STATUS CLASS */
type Job struct {
	Jobname string `bson:"jobname,omitempty" json:"jobname,omitempty"`
	Jobid   string `bson:"jobid,omitempty" json:"jobid,omitempty"`
	Owner   string `bson:"owner,omitempty" json:"owner,omitempty"`
	Status  string `bson:"status,omitempty" json:"status,omitempty"`
	Class   string `bson:"class,omitempty" json:"class,omitempty"`
}

func Dial(adr string, timeout time.Duration) (*Zftp, error) {
	var err error
	var z Zftp
	z.ServerConn, err = ftp.DialTimeout(adr, timeout)
	if err != nil {
		return nil, err
	} else {
		return &z, nil
	}
}

func (z *Zftp) SetSeqMode() (err error) {
	return z.Cmd("SITE FILETYPE=SEQ")
}

func (z *Zftp) SetJesMode() (err error) {
	return z.Cmd("SITE FILETYPE=JES")
}

func (z *Zftp) SetBinaryType() (err error) {
	return z.Cmd("TYPE I")
}

func (z *Zftp) SetAsciiType() (err error) {
	return z.Cmd("TYPE A")
}

func (z *Zftp) GetUnixFile(remote string) (r io.ReadCloser, err error) {
	return z.Retr(remote)
}

func (z *Zftp) PutUnixFile(r io.Reader, remote string) (err error) {
	conn, err := z.CmdDataConnFrom(0, "STOR %s", remote)
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, r)
	conn.Close()
	if err != nil {
		return err
	}

	_, _, err = z.GetConn().ReadResponse(ftp.StatusRequestedFileActionOK)
	return err
}

func (z *Zftp) DeleteUnixFile(remote string) (err error) {
	return z.Delete(remote)
}

func (z *Zftp) GetPsDataset(dataset string) (r io.ReadCloser, err error) {
	dsname := fmt.Sprintf("'%s'", dataset)
	return z.Retr(dsname)
}

func (z *Zftp) GetPdsDataset(dataset, dir string) (err error) {
	return nil
}

// The z/OS FTP server supports only the CRLF("\r\n") value for incoming ASCII data.
func (z *Zftp) SubmitJob(r io.Reader) (jobid string, err error) {
	z.generizeJesEnv()
	conn, err := z.CmdDataConnFrom(0, "STOR %s", "'ZFTP.X.Y.Z'")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(conn, r)
	conn.Close()
	if err != nil {
		return "", err
	}

	_, message, err := z.GetConn().ReadResponse(ftp.StatusRequestedFileActionOK)
	if err != nil {
		return "", err
	} else {
		// Get Job ID: J(OB|00)\d{5}
		// It is known to JES as J0013819
		re, _ := regexp.Compile(`It is known to JES as ([\w\d]{8})`)
		result := re.FindStringSubmatch(message)
		if result == nil {
			return "", fmt.Errorf("reponse text: %s", message)
		}
		// The number of fields in the resulting array always matches the number of groups plus one
		jobid := result[1]
		return jobid, nil
	}
}

func (z *Zftp) SubmitRemoteJob(dataset string) (jobid string, err error) {
	z.SetSeqMode()
	r, err := z.GetPsDataset(dataset)
	data, err := ioutil.ReadAll(r)
	z.SetJesMode() // will hang at HERE [TODO 2016-06-18]
	br := bytes.NewBuffer(data)
	fmt.Println("SubmitJob")
	return z.SubmitJob(br)
}

func (z *Zftp) PurgeJob(jobid string) (err error) {
	// JESJOBNAME=MEGA*, JESSTATUS=ALL and JESOWNER=MEGA
	return z.Delete(jobid)
}

func (z *Zftp) GetJoblogByID(jobid string) (joblog string, err error) {
	return "", nil
}

func (z *Zftp) GetJobStatusByID(jobid string) (j *Job, err error) {
	z.generizeJesEnv()
	conn, err := z.CmdDataConnFrom(0, "LIST %s", jobid)
	if err != nil {
		return nil, err
	}

	r := z.GetResponse(conn)
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	firstLine := scanner.Text()
	// The first line should be TiTLE
	//   JOBNAME  JOBID    OWNER    STATUS CLASS
	validTitle := regexp.MustCompile(`\s*JOBNAME\s+JOBID\s+OWNER\s+STATUS\s+CLASS`)
	if validTitle.MatchString(firstLine) == false {
		return nil, fmt.Errorf("Invalid list title: ", firstLine)
	}

	jobEntry := regexp.MustCompile(`\s*(\w+)\s+(\w+)\s+(\w+)\s+(\w+)\s+(\w+)`)
	scanner.Scan()
	line := scanner.Text()
	result := jobEntry.FindStringSubmatch(line)
	if result == nil {
		return nil, fmt.Errorf("Unmatched job entry: ", line)
	}
	//   JOBNAME  JOBID    OWNER    STATUS CLASS
	j = &Job{
		Jobname: result[1],
		Jobid:   result[2],
		Owner:   result[3],
		Status:  result[4],
		Class:   result[5],
	}

	return j, nil
}

func (z *Zftp) GetJobStatusByName(jobid string) (status string, err error) {
	return "", nil
}

func (z *Zftp) generizeJesEnv() error {
	var err error
	err = z.Cmd("SITE JESJOBNAME=*")
	if err != nil {
		return err
	}
	err = z.Cmd("SITE JESOWNER=*")
	if err != nil {
		return err
	}
	err = z.Cmd("SITE JESSTATUS=ALL")
	if err != nil {
		return err
	}
	return nil
}

func (z *Zftp) GetJobLog(jobid string) (r io.ReadCloser, err error) {
	z.generizeJesEnv()
	return z.Retr(jobid + ".x")
}

func (z *Zftp) Unix2Dos(text string) string {
	var buffer bytes.Buffer
	cr := false
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '\r', '\n':
			cr = true
		default:
			if cr {
				buffer.WriteString("\r\n")
			}
			cr = false
			buffer.WriteByte(text[i])
		}
	}
	if cr {
		buffer.WriteString("\r\n")
	}

	return buffer.String()
}
