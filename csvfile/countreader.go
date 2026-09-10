package csvfile

import (
	"io"
	"fmt"
)

type CountReader struct {
	R     io.Reader
	Count int64
	Limit int64
}

func (c *CountReader) Close() error {
	rc, ok := c.R.(io.ReadCloser)
	if !ok {
		return fmt.Errorf("stream is not a closable stream")
	}
	return rc.Close()
}

func (c *CountReader) Read(p []byte) (n int, err error) {
	n, err = c.R.Read(p)

	if c.Limit > 0 && c.Count+int64(n) > c.Limit {
		n = int(c.Limit - c.Count)
		if n == 0 {
			return 0, io.EOF
		}
	}
	c.Count += int64(n)
	return
}
