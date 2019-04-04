package darwind3

import (
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
	"time"
)

func (u *DarwinD3) ftpClient(f func(*ftp.ServerConn) error) error {
	log.Println("FTP: Connecting")
	con, err := ftp.DialTimeout(u.Config.D3.Ftp.Server, time.Minute)
	if err != nil {
		return err
	}

	log.Println("FTP: Login")
	if err := con.Login(u.Config.D3.Ftp.User, u.Config.D3.Ftp.Password); err != nil {
		log.Println("FTP:", err)
		con.Quit()
		return nil
	}

	if err := f(con); err != nil {
		log.Println("FTP: Quit", err)
		con.Quit()
		return err
	}

	log.Println("FTP: Quit")
	return con.Quit()
}

// FtpLs utility to log the files in a path
func FtpLs(con *ftp.ServerConn, path string) error {
	log.Println("FTP: LS", path)
	if entries, err := con.List(path); err != nil {
		return err
	} else {
		for i, e := range entries {
			log.Printf("FTP: %d %s %v %d %v\n", i, e.Name, e.Type, e.Size, e.Time)
		}
	}

	return nil
}

func FtpCp(con *ftp.ServerConn, srcPath, destPath string) error {
	log.Println("FTP: Retrieving", srcPath)
	src, err := con.Retr(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	c, err := io.Copy(dest, src)
	if err != nil {
		return err
	}

	log.Println("FTP: Written", c, "to", destPath)
	return nil
}
