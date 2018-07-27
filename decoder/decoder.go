package decoder

import (
	"errors"
	"log"

	"github.com/FgoDt/tsmux/common"
)

type Decoder struct {
	TSIO     common.TSMUXIO
	TSHeader *common.TSHeader
}

func Decode(d *Decoder) error {
	if d.TSIO == nil {
		return errors.New("no io reader")
	}

	e := ParseTSHeader(d)
	if e != nil {
		return errors.New("parse ts header error")
	}

	return nil
}

func ParseTSHeader(d *Decoder) error {

	data := make([]byte, 4)

	if data == nil {
		log.Println("no mem !")
		return errors.New("no mem")
	}

	e := d.TSIO.Read(data, 4)
	if e != nil {
		log.Println("io read error")
		return errors.New("io read error")
	}

	if data[0] != 0x47 {
		return errors.New("not ts header")
	}

	d.TSHeader = &common.TSHeader{}
	d.TSHeader.Sync_byte = 0x47

	temp := data[1]
	d.TSHeader.Transport_error_indieator = temp >> 7
	d.TSHeader.Payload_uint_start_indeicator = (temp >> 6) & 1
	d.TSHeader.Transport_priority = (temp >> 5) & 1
	d.TSHeader.PID = uint16(temp&0x1f)<<8 + uint16(data[2])
	temp = data[3]
	d.TSHeader.Transport_scrambling_control = (temp >> 6) & 3
	d.TSHeader.Adaptation_field_control = (temp >> 4) & 3
	d.TSHeader.Continuity_counter = temp & 0xf
	data = nil
	return nil
}
