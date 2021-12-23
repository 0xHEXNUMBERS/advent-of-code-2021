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

//This func is incorrect
//If the last byte contains 1-7 bits, the MSB will
//not be set correctly (0xxxxxxx where x == bits from the bitstream).
//Lukily, neither the test cases nor my specific input correctly
//tested for this behavior. Incidentally, Part 2's test cases also
//did not test for this, but my specific input did. This bug is
//fixed in 16/2/main.go
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
	return bSub
}

type Packet interface {
	VersionNumber() int
	SumOfVersions() int
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

		slb := b.streamLength
		subStream := createSubBitStream(b, int(subPacketsBitsLength))
		sla := b.streamLength
		if slb-sla != int(subPacketsBitsLength) {
			panic("SLB != SLA")
		}
		subPackets := make([]Packet, 0)
		for !subStream.IsFlushed() {
			log.Println("Reading Packet")
			log.Printf("SubBitStream: %#v", subStream)
			subPacket := readPacket(subStream)
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

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input = input[:len(input)-1] //Remove last

	fmt.Println(versionSum(string(input)))
}
