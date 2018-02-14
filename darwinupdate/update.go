// Package that handles FTP updates from the NRE FTP server
package darwinupdate

type DarwinUpdate struct {
  // The server name
  Server  string
  // The ftp user
  User    string
  // The ftp password for the NRE ftp server
  Pass    string
}
