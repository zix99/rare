package multiterm

import "os"

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
