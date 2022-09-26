package multiterm

import "os"

// Returns 'true' if output is being piped (Not char device)
func IsPipedOutput() bool {
	if fi, err := os.Stdout.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			return true
		}
	}
	return false
}

type BufferedTerm struct {
	*VirtualTerm
}

// NewBufferedTerm writes on Close
func NewBufferedTerm() *BufferedTerm {
	return &BufferedTerm{
		NewVirtualTerm(),
	}
}

func (s *BufferedTerm) Close() {
	s.WriteToOutput(os.Stdout)
	s.VirtualTerm.Close()
}
