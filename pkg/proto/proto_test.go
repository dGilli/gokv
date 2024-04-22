package proto

import (
	"fmt"
	"testing"
)

func TestProtocol(t *testing.T) {
    raw := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
    cmd, err := ParseCommand(raw)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(cmd)
}
