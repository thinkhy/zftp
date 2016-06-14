package zftp

import (
	"time"
	"io"
	//"bufio"
	ftp "./ftp"
	"strings"
	"fmt"
)

type Zftp struct {
	*ftp.ServerConn
}

func Dial(adr string, timeout time.Duration) (*Zftp, error) {
	var err error
	var z Zftp
	z.ServerConn, err = ftp.DialTimeout(adr, timeout*time.Second)
	if err != nil {
		return nil, err
	} else {
		return &z, nil
	}
}

/*
func (z *Zftp) Login(user, password string) (err error) {
	return nil
}

func (z *Zftp) Quit() (err error) {
	return nil
}
*/

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

	_, _, err = z.Connection().ReadResponse(ftp.StatusRequestedFileActionOK)
	return err
}

func (z *Zftp) GetPsDataset(remote, local string) (err error) {
	return nil
}

func (z *Zftp) GetPdsDataset(dataset, dir string) (err error) {
	return nil
}

// The z/OS FTP server supports only the CRLF("\r\n") value for incoming ASCII data.
func (z *Zftp) SubmitJob(r io.Reader) (jobid string, err error) {
	conn, err := z.CmdDataConnFrom(0, "STOR %s", "'ZFTP.X.Y.Z'")
	if err != nil {
		return "",err
	}

	_, err = io.Copy(conn, r)
	conn.Close()
	if err != nil {
		return "",err
	}

	_, message, err := z.Connection().ReadResponse(ftp.StatusRequestedFileActionOK)
	if err != nil {
		return "",err
	} else {
	        // Get Job ID: JOB\d{5}
		start := strings.Index(message, "JOB")
		if start == -1 {
			err = fmt.Errorf("Failed to submit job %s", message)
			return "",err
		}
		end := start + 3  // skim over "JOB"
		for end < len(message) && message[end] >= '0' && message[end] <= '9' {
		  end++
		}
		jobid := message[start:end]
		return jobid,nil
	}
	return "", nil
}

func (z *Zftp) SubmitRemoteJob(dataset string) (jobid string, err error) {
	return "", nil
}

func (z *Zftp) GetJoblogByID(jobid string) (joblog string, err error) {
	return "", nil
}

func (z *Zftp) GetJobStatusByID(jobid string) (status string, err error) {
	return "", nil
}

func (z *Zftp) GetJobStatusByName(jobid string) (status string, err error) {
	return "", nil
}


