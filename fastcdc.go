//go:generate go run _gen/gear.go -seed 1337 -output gear.go

package fastcdc

import (
	"bytes"
	"io"
	"sync"
)

// const minsize = 2 << 10

// const maskA = 0x0000d90303530000 // 13 ‘1’ bits

const avgsize = 8 << 10
const maxsize = 64 << 10

const maskS = 0x0003590703530000 // 15 ‘1’ bits
const maskL = 0x0000d90003530000 // 11 ‘1’ bits

var bufpool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, maxsize<<1))
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Do(r io.Reader, emit func(off int64, fp uint64, p []byte)) error {
	var lo, hi int64
	buf := bufpool.Get().(*bytes.Buffer)
	defer bufpool.Put(buf)

	for {
		if short := int64(maxsize - buf.Len()); short > 0 {
			_, err := buf.ReadFrom(io.LimitReader(r, short))
			if err != nil {
				return err
			}
		}

		if buf.Len() == 0 {
			return nil
		}

		b := buf.Bytes()
		i, fp := 0, uint64(0)

		for pivot := min(len(b), avgsize); i < pivot; i++ {
			if fp = (fp << 1) + gear[b[i]]; fp&maskS == 0 {
				goto emitchunk
			}
		}

		for max := min(len(b), maxsize); i < max; i++ {
			if fp = (fp << 1) + gear[b[i]]; fp&maskL == 0 {
				break
			}
		}

	emitchunk:
		lo, hi = hi, hi+int64(i)
		emit(lo, fp, buf.Next(i))
	}
}
