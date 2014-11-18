package event

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
)

type AggregateId [16]byte

var ZeroId AggregateId

func init() {
	gob.Register(AggregateId{})
}

// uuid generates a random UUID according to RFC 4122
func GenerateId() AggregateId {
	uuid := AggregateId{}
	n, err := io.ReadFull(rand.Reader, uuid[:])
	if n != len(uuid) || err != nil {
		panic("failed to create uuid")
	}
	uuid[8] = (uuid[8] | 0x80) & 0xBf
	uuid[6] = (uuid[6] | 0x40) & 0x4f
	return uuid
}

func (uuid AggregateId) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", (uuid)[0:4], (uuid)[4:6], (uuid)[6:8], (uuid)[8:10], (uuid)[10:])
}

func (uuid *AggregateId) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		if len(src) != 16 {
			return fmt.Errorf("invalid length %v", len(src))
		}
		copy((*uuid)[:], src)
	default:
		return errors.New("invalid type")
	}
	return nil
}
func (uuid AggregateId) Value() (driver.Value, error) {
	return []byte(uuid[:]), nil
}

func ParseId(s string) (id AggregateId, ok bool) {
	for i, x := range []int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		v, ok := xtob(s[x:])
		if !ok {
			return id, false
		}
		id[i] = v
	}
	return id, true
}

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = []byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// xtob converts the the first two hex bytes of x into a byte.
func xtob(x string) (byte, bool) {
	b1 := xvalues[x[0]]
	b2 := xvalues[x[1]]
	return (b1 << 4) | b2, b1 != 255 && b2 != 255
}
