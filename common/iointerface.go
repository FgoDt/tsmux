package common

type TSMUXIO interface {
	Read([]byte, uint64) error
	Write([]byte, uint64) error
}
