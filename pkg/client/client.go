package client

import (
	"bytes"
	"context"
	"net"

	"github.com/dGilli/gokv/pkg/proto"
	"github.com/tidwall/resp"
)

type (
    Client struct {
        addr string
    }
)

func New(addr string) *Client {
    return &Client{
        addr: addr,
    }
}

func (c *Client) Set(ctx context.Context, key, val string) error {
    conn, err := net.Dial("tcp", c.addr)
    if err != nil {
        return err
    }

    var buf bytes.Buffer
    wr := resp.NewWriter(&buf)
    wr.WriteArray([]resp.Value{
        resp.StringValue(proto.CommandSET),
        resp.StringValue(key),
        resp.StringValue(val),
    })
    _, err = conn.Write(buf.Bytes())
    return err
}
