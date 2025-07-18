package block

import (
	"encoding/binary"
	"io"
	"time"

	"go.uber.org/atomic"
)

const (
	TypeConnect = iota
	TypeDisconnect
	TypeData
	TypePing
	TypePong

	ShutdownRead = iota
	ShutdownWrite
	ShutdownBoth

	HeaderSize = 1 + 4 + 8 + 4 + 4
	DataSize   = 32*1024 - 13
	MaxSize    = HeaderSize + DataSize
)

type Block struct {
	Type         uint8  // 1 byte
	ConnectionID uint32 // 4 bytes
	TimeStamp    int64 // 8 bytes
	BlockID      uint32 // 4 bytes
	BlockLength  uint32 // 4 bytes
	BlockData    []byte
	packed       []byte
}

func (block *Block) Pack() []byte {
	if block.packed != nil {
		return block.packed
	}
	block.packed = make([]byte, HeaderSize+len(block.BlockData))
	block.packed[0] = block.Type
	binary.LittleEndian.PutUint32(block.packed[1:], block.ConnectionID)
	binary.LittleEndian.PutUint64(block.packed[5:], uint64(block.TimeStamp))
	binary.LittleEndian.PutUint32(block.packed[13:], block.BlockID)
	binary.LittleEndian.PutUint32(block.packed[17:], block.BlockLength)
	copy(block.packed[HeaderSize:], block.BlockData)
	return block.packed
}

func NewBlockFromReader(reader io.Reader) (*Block, error) {
	headerBuf := make([]byte, HeaderSize)
	block := Block{}
	_, err := io.ReadFull(reader, headerBuf)
	if err != nil {
		return nil, err
	}
	block.Type = headerBuf[0]
	block.ConnectionID = binary.LittleEndian.Uint32(headerBuf[1:])
	block.TimeStamp = int64(binary.LittleEndian.Uint64(headerBuf[5:]))
	block.BlockID = binary.LittleEndian.Uint32(headerBuf[13:])
	block.BlockLength = binary.LittleEndian.Uint32(headerBuf[17:])
	block.BlockData = make([]byte, block.BlockLength)
	if block.BlockLength > 0 {
		_, err = io.ReadFull(reader, block.BlockData)
		if err != nil {
			return nil, err
		}
	}
	return &block, nil
}

func NewConnectBlock(connectID uint32, blockID uint32, address string) Block {
	data := []byte(address)
	return Block{
		Type:         TypeConnect,
		ConnectionID: connectID,
		TimeStamp:    time.Now().UnixNano(),
		BlockID:      blockID,
		BlockLength:  uint32(len(data)),
		BlockData:    data,
	}
}

func newDataBlock(connectID uint32, blockID uint32, data []byte) Block {
	// We should copy data now
	blk := Block{
		Type:         TypeData,
		ConnectionID: connectID,
		TimeStamp:    time.Now().UnixNano(),
		BlockID:      blockID,
		BlockLength:  uint32(len(data)),
		BlockData:    data,
	}
	blk.Pack()
	return blk
}

func NewDataBlocks(connectID uint32, blockID *atomic.Uint32, data []byte) []Block {
	blocks := make([]Block, 0)
	for cursor := 0; cursor < len(data); {
		end := cursor + DataSize
		if len(data) < end {
			end = len(data)
		}
		blocks = append(blocks, newDataBlock(connectID, blockID.Inc()-1, data[cursor:end]))
		cursor = end
	}
	return blocks
}

func NewDisconnectBlock(connectID uint32, blockID uint32, shutdownType uint8) Block {
	return Block{
		Type:         TypeDisconnect,
		ConnectionID: connectID,
		TimeStamp:    time.Now().UnixNano(),
		BlockID:      blockID,
		BlockLength:  1,
		BlockData:    []byte{shutdownType},
	}
}

func NewPingBlock(connectID uint32, blockID uint32, lancenty uint64) Block {
	buf := make([]byte, 8) 
	binary.LittleEndian.PutUint64(buf, lancenty) 
	return Block{
		Type:         TypePing,
		ConnectionID: connectID,
		TimeStamp:    time.Now().UnixNano(),
		BlockID:      blockID,
		BlockLength:  8,
		BlockData:    buf,
	}
}

func NewPongBlock(connectID uint32, blockID uint32, timestamp uint64) Block {
	buf := make([]byte, 8) 
	binary.LittleEndian.PutUint64(buf, timestamp) 
	return Block{
		Type:         TypePong,
		ConnectionID: connectID,
		TimeStamp:    time.Now().UnixNano(),
		BlockID:      blockID,
		BlockLength:  8,
		BlockData:    buf,
	}
}
