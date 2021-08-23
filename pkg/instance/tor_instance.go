package instance

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TorStatus int

const (
	STARTING = iota
	RUNNING
	STOPPING
	STOPPED
)

type TorInstance struct {
	Port    uint16
	status  TorStatus
	process *os.Process
}

func NewTorInstance(port uint16) *TorInstance {
	return &TorInstance{Port: port, status: STARTING}
}

func (i *TorInstance) Run() {
	defer func() {
		i.status = STOPPED
	}()
	name := strconv.Itoa(int(i.Port))

	cmd := exec.Command("tor",
		"-f", fmt.Sprintf("resources/%s/%s.txt", name, name),
		"--CookieAuthentication", "0",
		"--SocksPort", name,
		"--DataDirectory", fmt.Sprintf("resources/%s", name),
		"--NewCircuitPeriod", "15",
		"--MaxCircuitDirtiness", "15",
		"--NumEntryGuards", "8",
		"--CircuitBuildTimeout", "5",
		"--ExitRelay", "0",
		"--RefuseUnknownExits", "0",
		"--ClientOnly", "1",
		"--StrictNodes", "1")
	stdoutIn, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		return
	}

	i.process = cmd.Process

	err = createDirsAndFiles(name)
	if err != nil {
		return
	}

	i.capture(stdoutIn)

	_ = cmd.Wait()

	deleteDirsAndFiles(name)
}

func (i *TorInstance) Stop() {
	_ = i.process.Kill()
}

func (i *TorInstance) WaitForStart() {
	for {
		if i.status != STARTING {
			break
		}
	}
}

func (i *TorInstance) WaitForStop() {
	for {
		if i.status == STOPPED {
			break
		}
	}
}

func (i *TorInstance) IsRunning() bool {
	return i.status == RUNNING
}

func createDirsAndFiles(name string) error {
	err := os.MkdirAll(
		fmt.Sprintf("resources/%s", name),
		0755)
	if err != nil {
		return err
	}

	file, fErr := os.OpenFile(
		fmt.Sprintf("resources/%s/%s.txt", name, name),
		os.O_CREATE | os.O_RDWR,
		0644)
	if fErr != nil {
		return fErr
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func deleteDirsAndFiles(name string) {
	_ = os.RemoveAll(fmt.Sprintf("resources/%s", name))
}

func (i *TorInstance) capture(r io.Reader) {
	buf := make([]byte, 1024)
	var out string
	for {
		n, err := r.Read(buf)
		if n > 0 {
			out = string(buf[:n])
			for _, v := range strings.Split(out, "\n") {
				// fmt.Println(v)
				if strings.HasSuffix(v, "Bootstrapped 100% (done): Done") {
					i.status = RUNNING
				}
				if strings.HasSuffix(v, "Address already in use. Is Tor already running?") {
					i.status = STOPPING
					continue
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				i.status = STOPPING
			}
			return
		}
	}
}
