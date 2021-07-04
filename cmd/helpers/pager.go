package helpers

import (
	"io"
	"os"
	"os/exec"
)

// TryWritePager will attempt to output everything in reader to a pager
// (Defaults to less, or whatever $PAGER is set to)
func TryWritePager(r io.Reader) error {
	cmd := getPagerCommand()

	w, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	io.Copy(w, r)
	w.Close()
	return cmd.Wait()
}

func getPagerCommand() *exec.Cmd {
	if pager := os.Getenv("PAGER"); pager != "" {
		return exec.Command(pager)
	}
	return exec.Command("less", "-R")
}
