package common

type TSHeader struct {
	Sync_byte                     uint8
	Transport_error_indieator     uint8
	Payload_uint_start_indeicator uint8
	Transport_priority            uint8
	PID                           uint16
	Transport_scrambling_control  uint8
	Adaptation_field_control      uint8
	Continuity_counter            uint8
}
