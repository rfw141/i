package cmdx

import (
	"bytes"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
	"os/exec"
	"strings"
)

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	log.Infof("cmd: %s", cmd.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errStr := strings.ReplaceAll(strings.TrimSpace(stderr.String()), "\r\n", "\n")
		for _, v := range strings.Split(errStr, "\n") {
			log.Logf(log.ERROR, &log.Option{NoSource: true}, "cmd failed with %s", v)
		}
		log.Errorf("%v", err)
		return err
	}
	if len(out.String()) > 0 {
		log.Infof("cmd out: %s", out.String())
	}
	return nil
}

func MustRun(name string, args ...string) {
	err := Run(name, args...)
	errorx.Must(err)
}
