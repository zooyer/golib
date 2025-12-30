package embed

import "testing"

func TestMustMalloc(t *testing.T) {
	var (
		n     int
		err   error
		buf   = make([]byte, 1024)
		data  = "Hello World!"
		block = MustMalloc(Size1KB)
	)

	if _, err = block.Write([]byte(data)); err != nil {
		t.Fatal(err)
	}

	if n, err = block.Read(buf); err != nil {
		t.Fatal(err)
	}

	if string(buf[:n]) != data {
		t.Fatal(string(buf[:n]), "!=", data)
	}
}
