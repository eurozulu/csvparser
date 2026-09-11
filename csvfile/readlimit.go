package csvfile

import "io"

type ReadLimit struct {
	Limit int64
	Used  int64
	in    io.ReadCloser
}

func (r *ReadLimit) Read(p []byte) (n int, err error) {
	if r.Used >= r.Limit {
		return 0, io.EOF
	}
	n, err = r.in.Read(p)
	if err != nil {
		return n, err
	}
	if r.Used+int64(n) > r.Limit {
		n = int(r.Limit - r.Used)
	}
	r.Used += int64(n)
	return
}

func (r ReadLimit) Close() error {
	return r.in.Close()
}
