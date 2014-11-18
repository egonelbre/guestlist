package event

import "testing"

func TestAggregateIdString(t *testing.T) {
	tests := map[string]AggregateId{
		"00000000-0000-0000-0000-000000000000": AggregateId{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"550e8400-e29b-41d4-a716-446655440000": AggregateId{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00},
	}

	for expected, uuid := range tests {
		if expected != uuid.String() {
			t.Fatalf("%v returned %s expected %s", uuid, uuid.String(), expected)
		}
	}
}
