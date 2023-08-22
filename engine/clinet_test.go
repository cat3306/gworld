package engine

import (
	"github.com/panjf2000/gnet/v2"
	"testing"
)

func TestClient(t *testing.T) {
	cli, err := gnet.NewClient(NewClientEvents())
	if err != nil {
		t.Fatal(err)
	}
	cnn, err := cli.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal(err)
	}
	cli.Start()
	t.Log(cnn)
	select {

	}
}
