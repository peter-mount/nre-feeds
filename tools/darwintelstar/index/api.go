package index

import (
	"github.com/peter-mount/go-kernel/v2/log"
	"os"
	"os/exec"
)

type Api struct {
	Cmd  *string `kernel:"flag,telstar-util,telstar-util path"`
	Host *string `kernel:"flag,telstar-host,Telstar hostname"`
	User *string `kernel:"flag,telstar-user,Telstar username"`
	Pass *string `kernel:"flag,telstar-pass,Telstar password"`
}

func (a *Api) Login() error {
	log.Println("Logging in")
	if err := a.Call("login", *a.Host, *a.User, *a.Pass); err != nil {
		return err
	}
	log.Println("Logged in")
	return nil
}

func (a *Api) Call(cmd string, args ...string) error {
	ary := []string{cmd}
	ary = append(ary, args...)
	exe := exec.Command(*a.Cmd, ary...)

	// TODO is this enough, thinking telstar-util might be too chatty on it's output
	if log.IsVerbose() {
		//exe.Stdout = os.Stdout
		exe.Stderr = os.Stderr
		//exe.Stdin = os.Stdin
	}

	return exe.Run()
}
