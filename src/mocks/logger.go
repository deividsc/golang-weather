package mocks

import (
	"bytes"
	"log"
)

var (
	buff       bytes.Buffer
	MockLogger = log.New(&buff, "", log.LstdFlags)
)
