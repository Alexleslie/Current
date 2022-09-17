package utils

import "Current/Grpc/codec"

const (
	GobType  codec.EncodeType = "application/gob"
	JsonType codec.EncodeType = "application/json"
)

const MagicNumber = 0x0529
