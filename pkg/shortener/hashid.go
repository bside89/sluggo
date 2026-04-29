package shortener

import (
	"github.com/bwmarrin/snowflake"
	hashids "github.com/speps/go-hashids/v2"
)

// base62Alphabet is the character set used to produce URL-safe short codes.
const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Shortener generates unique, collision-free short codes by encoding
// Snowflake IDs with the hashids algorithm (Base62 + secret salt).
type Shortener struct {
	node *snowflake.Node
	h    *hashids.HashID
}

// New creates a Shortener using the given secret key as salt and the given
// Snowflake node ID (0–1023) to guarantee uniqueness across distributed nodes.
func New(secretKey string, nodeID int64) (*Shortener, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}

	hd := hashids.NewData()
	hd.Salt = secretKey
	hd.Alphabet = base62Alphabet
	hd.MinLength = 6

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}

	return &Shortener{node: node, h: h}, nil
}

// Encode generates a new unique short code by encoding a fresh Snowflake ID.
// Because Snowflake IDs are monotonically unique, no database pre-check is
// required to guarantee the resulting hash is non-colliding.
func (s *Shortener) Encode() (string, error) {
	id := s.node.Generate().Int64()
	return s.h.EncodeInt64([]int64{id})
}
