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

	t := d.TSBuf.SeekFront(188)
	if !t {
		return errors.New("Seek")
	}

	for {
		//get sync
		ParseTSHeader(d)

		switch d.TSHeader.PID {
		case 0:
			log.Println("parse PAT")
			d.ParsePAT()
			break
		case 0x11:
			log.Println("parse SDT")
			d.ParseSDT()
			break
		}
	}

	return nil
}

func ParseTSHeader(d *Decoder) error {

	data := make([]byte, 4)

	if data == nil {
		log.Println("no mem !")
		return errors.New("no mem")
	}

	e := d.TSBuf.Read(data,4)
	//e := d.TSIO.Read(data, 4)
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
	if d.TSHeader.PayloadUintStartIndeicator == 1{
		e = d.TSBuf.Read(data,1)
		if e != nil{
			return e
		}
	//	len := data[0]

	}
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

	data := make([]uint8,11)

	err := d.TSBuf.Read(data,11)
	if err != nil{
		log.Println(err.Error())
		return  errors.New("read buf error")
	}

	sdt := new (common.SDTHeader)
	sdt.TableId = data[0]

	tmp := data[1]

	sdt.SectionSyntaxIndicator = tmp & 0x1
	sdt.SectionLength = uint16(tmp  & 0x0f)

	tmp = data[2]
	sdt.SectionLength += uint16(tmp)

	sdt.TransportStreamId = uint16(data[3] << 8)
	sdt.TransportStreamId += uint16(data[4])

	tmp = data[5]
	sdt.VersionNumber = uint8(tmp >> 1 & 0x1f)
	sdt.CurrentNextIndicator = uint8(tmp >> 7 & 0x1)
	//
	sdt.SectionNumber = data[6]
	sdt.LastSectionNumber = data[7]

	sdt.OriginalNetwordId = data[8]<<8
	sdt.OriginalNetwordId += data[9]

	data = make([]uint8,183-11)
	err = d.TSBuf.Read(data,183-11)
	if err != nil {
		return  err
	}
	index := 0
	for{
		sid := data[index] << 8
		index ++
		sid += data[index]
		index ++
		if(sid<0){
			break
		}
		val := data[index]
		index ++
		if val < 0 {
			break
		}
		descListLen := data[index]<<8
		index ++
		descListLen += data[index]
		index ++
		descListLen = uint8(int(descListLen) & 0xfff)
		if int(descListLen) > len(data) - index {
			break
		}
		for {
			tag := data[index]
			index ++
			descLen := data[index]
			index ++
			if descLen + 2 > descListLen{
				break;
			}

			switch tag {
			case 0x48:
				serverType := data[index]
				index ++
				if serverType < 0 {
					break;
				}

				len := data[index]
				index ++
				if len < 0 {
					break
				}
				providerName := string(data[index:index+int(len)])
				log.Println("providerName :",providerName)
				index += int(len)

				len = data[index]
				index ++
				name := string(data[index:index+int(len)])
				log.Println("programe name :",name)
				index += int(len)
				break
			default:
				break

			}
		}
	}



	return nil
}

func (d *Decoder) ParsePAT() error{
	data := make ([]uint8, 183)
	err := d.TSBuf.Read(data,183)
	if err != nil {
		return err
	}
	offset := 0
	pat := new(common.PAT)
	pat.TableID = data[offset]
	offset ++
	pat.SectionLen = uint16(data[offset] << 8)
	offset ++
	pat.SectionLen += uint16(data[offset])
	offset ++
	pat.SectionLen = pat.SectionLen & 0xfff
	pat.TransportStreamID = uint16(data[offset] << 8)
	offset ++
	pat.TransportStreamID += uint16(data[offset])
	offset ++
	pat.Version = data[offset]
	offset ++
	pat.Version = pat.Version >> 1 & 0x1f
	pat.SectionNumber = data[offset]
	offset ++
	pat.LastSectionNumber = data[offset]
	offset ++

	for{
		programNumber := data[offset] << 8
		offset ++
		programNumber += data[offset]
		offset ++
		if programNumber == 0 {
			networkID := int(data[offset]) << 8
			offset ++
			networkID += int(data[offset])
			offset ++
			networkID = networkID & 0x1fff
		}else {
			pmtID := int(data[offset])<<8
			offset ++
			pmtID += int(data[offset])
			offset ++
			pmtID = pmtID & 0x1fff
			log.Println("pmt id ",pmtID)
		}
	}

	return nil
}