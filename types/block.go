package types

import (
	"time"

	"google.golang.org/protobuf/proto"

	"authentication-chains/cipher"
)

func NewBlock(prevHash []byte, index uint64, dar *DeviceAuthenticationRequest) *Block {
	block := &Block{
		Hash:      nil,
		PrevHash:  prevHash,
		Index:     index,
		Dar:       dar,
		Timestamp: time.Now().Unix(),
	}

	data, err := proto.Marshal(block)
	if err != nil {
		panic(err)
	}

	block.Hash = cipher.Hash(data)

	return block
}

// Serialize serializes a block.
func (b *Block) Serialize() []byte {
	data, err := proto.Marshal(b)
	if err != nil {
		panic("serialize block failed: " + err.Error())
	}

	return data
}

// DeserializeBlock deserializes a block.
func DeserializeBlock(data []byte) *Block {
	block := &Block{}

	err := proto.Unmarshal(data, block)
	if err != nil {
		panic("deserialize block failed: " + err.Error())
	}

	return block
}
