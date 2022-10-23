package jobs

import (
	"bytes"
	"os"
	"testing"
)

func TestCompressFile(t *testing.T) {
	zippedFileName, zippedFileBuffer, _ := compressFile("test.txt", bytes.NewBuffer([]byte("Test string")))
	f, _ := os.Create(zippedFileName)
	f.Write(zippedFileBuffer.Bytes())
	f.Close()
}
