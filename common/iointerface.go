package common

type TSMUXIO interface {
	Read([]byte, int) error
	Write([]byte, int) error
}
