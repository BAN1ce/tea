package mqtt

import (
	"bufio"
	"errors"
	"tea/src/mqtt/handle"
)

/**
 * CONNECT Packet.
 */
const CMD_CONNECT = 1;

/**
 * CONNACK
 */
const CMD_CONNACK = 2;

/**
 * PUBLISH
 */
const CMD_PUBLISH = 3;

/**
 * PUBACK
 */
const CMD_PUBACK = 4;

/**
 * PUBREC
 */
const CMD_PUBREC = 5;

/**
 * PUBREL
 */
const CMD_PUBREL = 6;

/**
 * PUBCOMP
 */
const CMD_PUBCOMP = 7;

/**
 * SUBSCRIBE
 */
const CMD_SUBSCRIBE = 8;

/**
 * SUBACK
 */
const CMD_SUBACK = 9;

/**
 * UNSUBSCRIBE
 */
const CMD_UNSUBSCRIBE = 10;

/**
 * UNSUBACK
 */
const CMD_UNSUBACK = 11;

/**
 * PINGREQ
 */
const CMD_PINGREQ = 12;

/**
 * PINGRESP
 */
const CMD_PINGRESP = 13;

/**
 * DISCONNECT
 */
const CMD_DISCONNECT = 14;

/**
mqtt 包解析
*/
func Input() bufio.SplitFunc {

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		if atEOF {
			return len(data), data[:len(data)], errors.New("EOF")
		}

		if len(data) >= 2 {
			headLength := 1
			multiplier := 1
			bodyLength := 0
			for i := 1; i < len(data); i++ {
				bodyLength += int(data[i]&127) * multiplier
				headLength++
				if data[i]&128 == 1 {
					multiplier *= 128;
				} else {
					break
				}
			}
			sum := headLength + bodyLength
			return sum, data[0:sum], nil
		}

		return 0, nil, nil
	}
}

func Decode(data []byte) *handle.Pack {
	pack := handle.NewPack()

	pack.Cmd = int(data[0] >> 4)
	fixHeadLength := 1
	multiplier := 1
	bodyLength := 0
	for i := 1; i < len(data); i++ {
		bodyLength += int(data[i]&127) * multiplier
		fixHeadLength++
		if data[i]&128 == 1 {
			multiplier *= 128;
		} else {
			break
		}
	}
	pack.FixedHeader = data[0:fixHeadLength]
	pack.FixHeaderLength = fixHeadLength
	pack.Data = data

	return pack
}

func Encode(pack handle.Pack) {

}
