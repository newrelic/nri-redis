package helpers

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// ExecInContainer executes the given command inside the specified container. It returns three values:
// 1st - Standard Output
// 2nd - Standard Error
// 3rd - Runtime error, if any
func ExecInContainer(container string, command []string, envVars ...string) (string, string, error) {
	cmdLine := make([]string, 0, 3+len(command))
	cmdLine = append(cmdLine, "exec", "-i")

	for _, envVar := range envVars {
		cmdLine = append(cmdLine, "-e", envVar)
	}

	cmdLine = append(cmdLine, container)
	cmdLine = append(cmdLine, command...)

	log.Printf("executing: docker %s", strings.Join(cmdLine, " "))

	fmt.Println(cmdLine)
	cmd := exec.Command("docker", cmdLine...)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()

	if err != nil {
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}
