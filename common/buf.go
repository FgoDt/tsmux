package common

import "errors"

func (buf *TSMUXBuf) Init(io TSMUXIO, maxSize uint64) error {
	buf.MaxSize = maxSize
	buf.Len = 0
	buf.Pos = 0
	buf.buf = make([]uint8, maxSize)
	buf.io = io
	return nil
}

func (b *TSMUXBuf) Read(data []uint8, len uint64) error {
	if len > 204 {
		return errors.New("can not read data more than 204 byte")
	}
	canRead := len
	if b.Len-b.Pos < len {
		canRead = b.Len - b.Pos
	}
	copy(data, b.buf[b.Pos:b.Pos+canRead])
	b.Pos += canRead

	if canRead != len {
		err := b.io.Read(b.buf, b.MaxSize)
		if err != nil {
			return err
		}
		b.Pos = 0
		b.Len = b.MaxSize
	}

	less := len - canRead
	if less > 0 {
		copy(data[canRead:], b.buf[b.Pos:less])
		b.Pos += less
	}
	return nil
}

func (b *TSMUXBuf) Write(data []uint8, len uint64) error {
	return nil
}

func (b *TSMUXBuf) Seek(pos uint64) bool {
	if pos > b.Len || pos < 0 {
		return false
	}
	b.Pos = pos
	return true
}
