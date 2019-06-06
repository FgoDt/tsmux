package decoder

import (
	"errors"
	"log"

	"github.com/FgoDt/tsmux/common"
)

type Decoder struct {
	TSIO     common.TSMUXIO
	TSHeader *common.TSHeader
	TSAdap   *common.TSAdaptationField
	SDTMap   *common.SDTHeader
	TSBuf    *common.TSMUXBuf
}

func (d *Decoder) Init() error {
	if d.TSIO == nil {
		return errors.New("no io reader")
	}

	d.TSBuf = &common.TSMUXBuf{}

	err := d.TSBuf.Init(d.TSIO, 1400)
	if err != nil {
		log.Println(err.Error())
		return errors.New("init tsio error ")
	}
	return nil
}

func (d *Decoder) Run() error {
	buf := make([]uint8, 188)
	e := d.TSBuf.Read(buf, uint64(len(buf)))
	if e != nil {
		return e
	}
	i := 0
	for i < 188 {
		if buf[i] == 0x47 {
			break
		}
		i++
	}
	if i == 188 {
		return errors.New("no sync find")
	}

	t := d.TSBuf.Seek(uint64(i))
	if !t {
		return errors.New("Seek")
	}

	//get sync
	ParseTSHeader(d)

	switch d.TSHeader.PID {
	case 0:
		log.Println("parse PAT")
		break
	case 0x11:
		log.Println("parse SDT")
		d.ParseSDT()
		break
	}

	return nil
}

func ParseTSHeader(d *Decoder) error {

	data := make([]byte, 4)
	data[0] = 0x47

	if data == nil {
		log.Println("no mem !")
		return errors.New("no mem")
	}

	e := d.TSIO.Read(data[1:], 3)
	if e != nil {
		log.Println("io read error")
		return errors.New("io read error")
	}

	if data[0] != 0x47 {
		return errors.New("not ts header")
	}

	d.TSHeader = &common.TSHeader{}
	d.TSHeader.SyncByte = 0x47

	temp := data[1]
	d.TSHeader.TransportErrorIndieator = temp >> 7
	d.TSHeader.PayloadUintStartIndeicator = (temp >> 6) & 1
	d.TSHeader.TransportPriority = (temp >> 5) & 1
	d.TSHeader.PID = uint16(temp&0x1f)<<8 + uint16(data[2])
	temp = data[3]
	d.TSHeader.TransportScramblingControl = (temp >> 6) & 3
	d.TSHeader.AdaptationFieldControl = (temp >> 4) & 3
	d.TSHeader.ContinuityCounter = temp & 0xf
	data = nil
	return nil
}

func ParseAdaptation(d *Decoder) error {
	return nil
}

func (d *Decoder) ParseSDT() error {
	if d.TSHeader.PID != 0x11 {
		return errors.New("not SDT pkt")
	}

	return nil
}
