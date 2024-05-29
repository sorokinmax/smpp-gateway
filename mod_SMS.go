package main

import (
	"maps"
)

const UDHIE_HEADER_LENGTH byte = 0x05
const UDHIE_IDENTIFIER_SAR byte = 0x00
const UDHIE_SAR_LENGTH byte = 0x03
const MAX_MULTIPART_MSG_SEGMENT_SIZE_UCS2 byte = 134

type Message struct {
	message    []byte
	partsCount byte
}

var Messages map[byte]map[byte]Message

func GetMessage(message string) string {
	messageBytes := []byte(message)

	if messageBytes[0] == UDHIE_HEADER_LENGTH && messageBytes[1] == UDHIE_IDENTIFIER_SAR && messageBytes[2] == UDHIE_SAR_LENGTH {
		id, partsCount, part := messageBytes[3], messageBytes[4], messageBytes[5]
		msg := Message{partsCount: partsCount}

		msg.message = messageBytes[6:]

		partOfMsg := map[byte]Message{part: msg}
		if len(Messages[id]) == 0 {
			Messages[id] = partOfMsg
		} else {
			maps.Copy(Messages[id], partOfMsg)
		}

		if Messages[id][part].partsCount == byte(len(Messages[id])) {
			var msg []byte
			for i := 1; i <= int(Messages[id][part].partsCount); i++ {
				msg = append(msg, Messages[id][byte(i)].message...)
			}
			return string(msg)
		}

	} else {
		return message
	}

	return ""
}
