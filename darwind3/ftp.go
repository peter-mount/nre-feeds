package darwind3

import (
	"github.com/jlaffaye/ftp"
	"log"
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
	log.Println("LS", path)
	if entries, err := con.List(path); err != nil {
		return err
	} else {
		for i, e := range entries {
			log.Printf("%d %s %v %d %v\n", i, e.Name, e.Type, e.Size, e.Time)
		}
	}

	return nil
}
