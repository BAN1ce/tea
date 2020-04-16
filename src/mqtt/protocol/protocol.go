package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"tea/src/manage"
	"tea/src/mqtt/response"
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

type Pack struct {
	Data            []byte
	PackLength      int
	FixedHeader     []byte
	FixHeaderLength int
	BodyLength      int
	Cmd             int
}

func NewPack() *Pack {

	pack := new(Pack)

	return pack
}

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

func Decode(data []byte) *Pack {
	pack := NewPack()

	pack.Cmd = int(data[0] >> 4)

	fmt.Println("cmd", pack.Cmd)
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
	pack.BodyLength = bodyLength
	pack.FixedHeader = data[0:fixHeadLength]
	pack.FixHeaderLength = fixHeadLength
	pack.Data = data
	pack.PackLength = len(data)

	return pack
}

func Encode(response response.Response, client *manage.Client) {

	variableHeaderLength := 0
	payloadLength := 0
	variableHeader, ok := response.GetVariableHeader()
	if ok {
		variableHeaderLength = len(variableHeader)
	}
	payload, ok := response.GetPayload()
	if ok == true {
		payloadLength = len(payload)
	}
	totalLength := variableHeaderLength + payloadLength

	responseData := make([]byte, 0)

	responseData = append(responseData, response.GetFixedHeaderWithoutLength())

	for {
		digit := totalLength % 128
		if digit > 0 {
			if totalLength/128 > 0 {
				digit = digit | 0x80
			}
			totalLength /= 128
			responseData = append(responseData, uint8(digit))
		} else {
			break
		}
	}
	if variableHeaderLength != 0 {
		responseData = append(responseData, variableHeader...)
	}
	if payloadLength != 0 {
		responseData = append(responseData, payload...)
	}

	fmt.Println("responseData", responseData)

	client.Write(responseData)
}
