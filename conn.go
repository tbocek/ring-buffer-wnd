package tomtp

import (
	"crypto/ecdh"
	"fmt"
	"net"
	"sync"
)

type ConnectionState uint8

const (
	ConnectionStarting ConnectionState = iota
	ConnectionEnding
	ConnectionEnded
)

type Connection struct {
	remoteAddr      net.Addr
	streams         map[uint32]*Stream
	listener        *Listener
	pubKeyIdRcv     *ecdh.PublicKey
	prvKeyEpSnd     *ecdh.PrivateKey
	pubKeyEpRcv     *ecdh.PublicKey
	sharedSecret    []byte
	rtoMillis       uint64
	nextSleepMillis uint64
	rbSnd           *RingBufferSnd[[]byte] // Send buffer for outgoing data, handles the global sn
	bytesWritten    uint64
	mtu             int
	sender          bool
	firstPaket      bool
	mu              sync.Mutex
	state           ConnectionState
}

func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, stream := range c.streams {
		//pick first stream, send close flag to close all streams
		stream.CloseAll()
		break
	}

	clear(c.streams)
	return nil
}

func (c *Connection) NewStreamSnd(streamId uint32) (*Stream, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.streams[streamId]; !ok {
		s := &Stream{
			streamId:     streamId,
			streamSnNext: 0,
			state:        StreamStarting,
			conn:         c,
			rbRcv:        NewRingBufferRcv[[]byte](maxRingBuffer, maxRingBuffer),
			writeBuffer:  []byte{},
			mu:           sync.Mutex{},
		}
		c.streams[streamId] = s
		return s, nil
	} else {
		return nil, fmt.Errorf("stream %x already exists", streamId)
	}
}

func (c *Connection) GetOrNewStreamRcv(streamId uint32) (*Stream, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stream, ok := c.streams[streamId]; !ok {
		s := &Stream{
			streamId:     streamId,
			streamSnNext: 0,
			state:        StreamStarting,
			conn:         c,
			rbRcv:        NewRingBufferRcv[[]byte](1, maxRingBuffer),
			writeBuffer:  []byte{},
			mu:           sync.Mutex{},
		}
		c.streams[streamId] = s
		return s, true
	} else {
		return stream, false
	}
}
