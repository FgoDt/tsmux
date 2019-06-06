package common

type TSHeader struct {
	SyncByte                   uint8
	TransportErrorIndieator    uint8
	PayloadUintStartIndeicator uint8
	TransportPriority          uint8
	PID                        uint16
	TransportScramblingControl uint8
	AdaptationFieldControl     uint8
	ContinuityCounter          uint8
}

type TSAdaptationField struct {
	AdaptationFieldLen                     uint8
	DiscontinuityIndicator                 uint8
	RandomAccessIndicator                  uint8
	ElementaryStreamPriorityIndicator      uint8
	PCRFlag                                uint8
	OPCRFlag                               uint8
	SplicingPointFlag                      uint8
	TransportPrivateDataFlag               uint8
	AdaptationFieldExtensionFlag           uint8
	ProgramClockReferenceBase              uint64 //pcr_flag == 1
	ProgramClockReferenceExtension         uint16
	OriginalProgramClockReferenceBase      uint64 //opcr flag == 1
	OriginalProgramClockReferenceExtension uint16
	SpliceCountdown                        uint8 //splicing point flag == 1
	TransportPrivateDataLen                uint8 //transport private data flag
	TransportPrivateData                   []uint8
	AdaptationFieldExtensionLen            uint8 //adaptation field extension flag
	LtwFlag                                uint8
	PiecewiseRateFlag                      uint8
	SeamlessSpliceFlag                     uint8
	LtwValidFlag                           uint8 //ltw flag
	LtwOffset                              uint16
	PiecewiseRate                          uint32 //piecewise rate flag
	SpliceType                             uint8  //seamless splice flag
	DTSNextAU                              []int8
}

type SDTProgram struct {
	ServerId                uint16
	EITScheduleFlag         uint8
	EITPresentFollowingFlag uint8
	RunningStatus           uint8 // 1  2  3 4
	FreeCAMode              uint8
	DescLoopLen             uint16
}

type SDTHeader struct {
	TableId                uint8
	SectionSyntaxIndicator uint8
	SectionLength          uint16
	TransportStreamId      uint16
	VersionNumber          uint8
	CurrentNextIndicator   uint8
	OriginalNetwordId      uint8
	programs               []SDTProgram
}

type TSMUXBuf struct {
	Len     uint64
	Pos     uint64
	MaxSize uint64
	buf     []uint8
	io      TSMUXIO
}
