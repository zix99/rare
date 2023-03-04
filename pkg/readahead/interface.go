package readahead

type OnScannerError func(error)

type Scanner interface {
	Scan() bool    // Scans for next string. True if exist, false if eof
	Bytes() []byte // Returns the result of the Scan()
	OnError(f OnScannerError)
}
