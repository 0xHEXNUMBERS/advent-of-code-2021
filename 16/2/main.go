package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type BitStream struct {
	buffer       []byte
	bitsRead     int
	streamLength int
	err          error
}

func (b *BitStream) ReadBit() byte {
	if b.err != nil {
		return 0
	}
	if b.streamLength == 0 {
		b.err = errors.New("EOS")
		return 0
	}

	if b.bitsRead == 8 {
		if len(b.buffer) <= 1 {
			b.err = errors.New("Reached EOB before supposed EOS")
			return 0
		}

		b.buffer = b.buffer[1:]
		b.bitsRead = 0
	}
	var ret byte = b.buffer[0] & 0x80
	ret >>= 7 //Make bit LSB

	b.buffer[0] <<= 1
	b.bitsRead++
	b.streamLength--
	return ret
}

func (b *BitStream) ReadBits(bits int) uint64 {
	if b.err != nil {
		return 0
	}

	var ret uint64
	if bits > 64 {
		b.err = errors.New("Cannot read more than 64 bits at a time")
		return 0
	}

	for i := 0; i < bits; i++ {
		bit := b.ReadBit()
		ret <<= 1
		ret |= uint64(bit)
	}
	log.Printf("READ %d BITS: %b", bits, ret)
	return ret
}

func (b *BitStream) IsFlushed() bool {
	return b.streamLength == 0
}

//hexToBin is valid for any rune 0-9 or A-F
func hexToBin(b byte) byte {
	if b >= '0' && b <= '9' {
		return b - '0'
	}
	return (b - 'A') + 0xA
}

func convertHexToBitStream(s string) *BitStream {
	b := &BitStream{}
	b.buffer = make([]byte, len(s)/2)
	b.streamLength = len(s) * 4 //hex digit represents 4 bits

	for i := 0; i < len(s)/2; i++ {
		hi := hexToBin(s[2*i])
		lo := hexToBin(s[2*i+1])
		b.buffer[i] = (hi << 4) | lo
	}
	return b
}

func createSubBitStream(bBase *BitStream, length int) *BitStream {
	numBytes := length / 8
	if length%8 > 0 {
		numBytes++
	}

	bSub := &BitStream{}
	bSub.buffer = make([]byte, numBytes)
	bSub.streamLength = length

	byteIndex := 0
	bitIndex := 0
	for i := 0; i < length; i++ {
		bit := bBase.ReadBit()
		bSub.buffer[byteIndex] <<= 1
		bSub.buffer[byteIndex] |= bit

		bitIndex++
		if bitIndex >= 8 {
			bitIndex = 0
			byteIndex++
		}
	}
	//Shift remainder bits to MSB
	if bitIndex != 0 {
		bSub.buffer[byteIndex] <<= (8 - bitIndex)
	}

	return bSub
}

type Packet interface {
	VersionNumber() int
	SumOfVersions() int
	PacketValue() uint64
}

type PacketHeader struct {
	Version int
	Type    int
}

func (p PacketHeader) VersionNumber() int {
	return p.Version
}

type LiteralPacket struct {
	PacketHeader
	Value uint64
}

func (l LiteralPacket) SumOfVersions() int {
	return l.Version
}

func (l LiteralPacket) PacketValue() uint64 {
	return l.Value
}

type OperatorPacket struct {
	PacketHeader
	SubPackets []Packet
}

func (o OperatorPacket) SumOfVersions() int {
	sum := o.Version
	for _, sub := range o.SubPackets {
		sum += sub.SumOfVersions()
	}
	return sum
}

func (o OperatorPacket) PacketValue() uint64 {
	const SUM = 0
	const PRODUCT = 1
	const MIN = 2
	const MAX = 3
	const GT = 5
	const LT = 6
	const EQ = 7
	switch o.Type {
	case SUM:
		var sum uint64
		for _, sub := range o.SubPackets {
			sum += sub.PacketValue()
		}
		return sum
	case PRODUCT:
		var acc uint64 = 1
		for _, sub := range o.SubPackets {
			acc *= sub.PacketValue()
		}
		return acc
	case MIN:
		min := o.SubPackets[0].PacketValue()
		for i := 1; i < len(o.SubPackets); i++ {
			val := o.SubPackets[i].PacketValue()
			if min > val {
				min = val
			}
		}
		return min
	case MAX:
		max := o.SubPackets[0].PacketValue()
		for i := 1; i < len(o.SubPackets); i++ {
			val := o.SubPackets[i].PacketValue()
			if max < val {
				max = val
			}
		}
		return max
	case GT:
		first := o.SubPackets[0].PacketValue()
		second := o.SubPackets[1].PacketValue()
		if first > second {
			return 1
		}
		return 0
	case LT:
		first := o.SubPackets[0].PacketValue()
		second := o.SubPackets[1].PacketValue()
		if first < second {
			return 1
		}
		return 0
	case EQ:
		first := o.SubPackets[0].PacketValue()
		second := o.SubPackets[1].PacketValue()
		if first == second {
			return 1
		}
		return 0
	}
	return 0
}

func readLiteralPacket(b *BitStream, header PacketHeader) LiteralPacket {
	var value uint64
	log.Printf("LITERAL START: %#v", b)
	for {
		bits := b.ReadBits(5)
		value <<= 4
		value |= bits & 0xF
		log.Printf("BITS READ: %X", bits)
		log.Printf("CURRENT LITERAL VALUE: %d", value)
		if bits&0x10 == 0 {
			break
		}
	}
	log.Printf("REST OF BITSTREAM: %#v", b)
	log.Printf("LITERAL VALUE: %d", value)
	return LiteralPacket{header, value}
}

func readOperatorPacket(b *BitStream, header PacketHeader) OperatorPacket {
	const NUM_BITS = 0
	const NUM_PACKETS = 1
	lengthTypeID := b.ReadBit()

	op := OperatorPacket{header, nil}
	if lengthTypeID == NUM_BITS {
		subPacketsBitsLength := b.ReadBits(15)
		log.Printf("NUM_BITS: %d", subPacketsBitsLength)
		log.Printf("BITSTREAM: %#v", b)

		subStream := createSubBitStream(b, int(subPacketsBitsLength))
		subPackets := make([]Packet, 0)
		for !subStream.IsFlushed() {
			log.Println("Reading Packet")
			log.Printf("SubBitStream BEFORE: %#v", subStream)
			subPacket := readPacket(subStream)
			log.Printf("SubBitStream AFTER: %#v", subStream)
			subPackets = append(subPackets, subPacket)
		}
		log.Println("SubPackets Read")
		log.Printf("Rest of BitStream: %#v", b)
		op.SubPackets = subPackets
	} else if lengthTypeID == NUM_PACKETS {
		numSubPackets := b.ReadBits(11)

		subPackets := make([]Packet, numSubPackets)
		log.Printf("NUMBER OF SUB PAKCETS: %d", numSubPackets)
		for i := 0; i < int(numSubPackets); i++ {
			log.Println("Reading Packet (1)")
			subPackets[i] = readPacket(b)
		}
		op.SubPackets = subPackets
	}
	return op
}

func readPacket(b *BitStream) Packet {
	const LITERAL_TYPE = 4

	header := PacketHeader{}
	header.Version = int(b.ReadBits(3))
	header.Type = int(b.ReadBits(3))

	log.Printf("Read Header: %#v", header)

	var packet Packet
	if header.Type == LITERAL_TYPE {
		log.Println("LITERAL")
		packet = readLiteralPacket(b, header)
	} else {
		log.Println("OPERATOR")
		packet = readOperatorPacket(b, header)
	}
	return packet
}

func versionSum(s string) int {
	bitStream := convertHexToBitStream(s)
	log.Printf("Initial Bit Stream: %#v", bitStream)
	packet := readPacket(bitStream)
	return packet.SumOfVersions()
}

func packetValue(s string) uint64 {
	bitStream := convertHexToBitStream(s)
	log.Printf("Initial Bit Stream: %#v", bitStream)
	packet := readPacket(bitStream)
	return packet.PacketValue()
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	fmt.Println(packetValue(string(input)))
}
