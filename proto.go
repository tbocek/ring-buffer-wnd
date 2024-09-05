package tomtp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

const protoHeaderSize = 24

type Payload struct {
	StreamId   uint32
	CloseFlag  bool
	RcvWndSize uint32
	AckStartSn uint32
	RleAck     uint64
	Sn         uint32
	Data       []byte
}

func EncodePayload(
	streamId uint32,
	closeFlag bool,
	rcvWndSize uint32,
	ackStartSn uint32,
	rleAck uint64,
	sn uint32,
	data []byte,
	w io.Writer) (n int, err error) {

	buf := new(bytes.Buffer)

	// STREAM_ID (32-bit)
	if err := binary.Write(buf, binary.BigEndian, streamId); err != nil {
		return 0, err
	}

	// Combine the close flag with the RCV_WND_SIZE
	if closeFlag {
		rcvWndSize |= 0x80000000 // Set the highest bit to 1
	}
	// RCV_WND_SIZE (31 bits) + STREAM_CLOSE_FLAG (1 bit)
	if err := binary.Write(buf, binary.BigEndian, rcvWndSize); err != nil {
		return 0, err
	}

	// ACK_START_SN (32-bit)
	if err := binary.Write(buf, binary.BigEndian, ackStartSn); err != nil {
		return 0, err
	}

	// RLE_ACK (64-bit)
	if err := binary.Write(buf, binary.BigEndian, rleAck); err != nil {
		return 0, err
	}

	// SEQ_NR (32-bit)
	if err := binary.Write(buf, binary.BigEndian, sn); err != nil {
		return 0, err
	}

	// Write DATA if present
	if len(data) > 0 {
		if _, err := buf.Write(data); err != nil {
			return 0, err
		}
	}

	n, err = w.Write(buf.Bytes())
	return n, err
}

func DecodePayload(buf *bytes.Buffer, n int) (payload *Payload, err error) {
	payload = &Payload{}
	bytesRead := 0

	// Helper function to read bytes and keep track of the count
	readBytes := func(num int) ([]byte, error) {
		if bytesRead+num > n {
			return nil, errors.New("Attempted to read " + strconv.Itoa(num) + " bytes when only " + strconv.Itoa(n-bytesRead) + " remaining")
		}
		b := make([]byte, num)
		_, err := io.ReadFull(buf, b)
		if err != nil {
			return nil, err
		}
		bytesRead += num
		return b, nil
	}

	// STREAM_ID (32-bit)
	streamIdBytes, err := readBytes(4)
	if err != nil {
		return nil, err
	}
	payload.StreamId = binary.BigEndian.Uint32(streamIdBytes)

	// RCV_WND_SIZE + STREAM_CLOSE_FLAG (32-bit)
	rcvWndSizeBytes, err := readBytes(4)
	if err != nil {
		return nil, err
	}
	rcvWndSize := binary.BigEndian.Uint32(rcvWndSizeBytes)
	payload.CloseFlag = (rcvWndSize & 0x80000000) != 0 // Extract the STREAM_CLOSE_FLAG
	payload.RcvWndSize = rcvWndSize & 0x7FFFFFFF       // Mask out the close flag to get the actual RCV_WND_SIZE

	// ACK_START_SN (32-bit)
	ackStartSnBytes, err := readBytes(4)
	if err != nil {
		return nil, err
	}
	payload.AckStartSn = binary.BigEndian.Uint32(ackStartSnBytes)

	// RLE_ACK (64-bit)
	rleAckBytes, err := readBytes(8)
	if err != nil {
		return nil, err
	}
	payload.RleAck = binary.BigEndian.Uint64(rleAckBytes)

	// SEQ_NR (32-bit)
	seqNrBytes, err := readBytes(4)
	if err != nil {
		return nil, err
	}
	payload.Sn = binary.BigEndian.Uint32(seqNrBytes)

	// Read the remaining data
	remainingBytes := n - bytesRead
	if remainingBytes > 0 {
		payload.Data = make([]byte, remainingBytes)
		_, err = io.ReadFull(buf, payload.Data)
		if err != nil {
			return nil, err
		}
	}

	return payload, nil
}
