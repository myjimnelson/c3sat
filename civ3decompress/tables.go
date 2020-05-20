package civ3decompress

type lengthCode struct {
	value, extraBits int
}

type bitKey struct {
	keyBitLength int
	key          uint
}

var lengthLookup = map[bitKey]lengthCode{
	{3, 0b101}:     {2, 0},
	{2, 0b11}:      {3, 0},
	{3, 0b100}:     {4, 0},
	{3, 0b011}:     {5, 0},
	{4, 0b0101}:    {6, 0},
	{4, 0b0100}:    {7, 0},
	{4, 0b0011}:    {8, 0},
	{5, 0b00101}:   {9, 0},
	{5, 0b00100}:   {10, 1},
	{5, 0b00011}:   {12, 2},
	{5, 0b00010}:   {16, 3},
	{6, 0b000011}:  {24, 4},
	{6, 0b000010}:  {40, 5},
	{6, 0b000001}:  {72, 6},
	{7, 0b0000001}: {136, 7},
	{7, 0b0000000}: {264, 8},
}

var offsetLookup = map[bitKey]int{
	{2, 0b11}:       0x00,
	{4, 0b1011}:     0x01,
	{4, 0b1010}:     0x02,
	{5, 0b10011}:    0x03,
	{5, 0b10010}:    0x04,
	{5, 0b10001}:    0x05,
	{5, 0b10000}:    0x06,
	{6, 0b011111}:   0x07,
	{6, 0b011110}:   0x08,
	{6, 0b011101}:   0x09,
	{6, 0b011100}:   0x0a,
	{6, 0b011011}:   0x0b,
	{6, 0b011010}:   0x0c,
	{6, 0b011001}:   0x0d,
	{6, 0b011000}:   0x0e,
	{6, 0b010111}:   0x0f,
	{6, 0b010110}:   0x10,
	{6, 0b010101}:   0x11,
	{6, 0b010100}:   0x12,
	{6, 0b010011}:   0x13,
	{6, 0b010010}:   0x14,
	{6, 0b010001}:   0x15,
	{7, 0b0100001}:  0x16,
	{7, 0b0100000}:  0x17,
	{7, 0b0011111}:  0x18,
	{7, 0b0011110}:  0x19,
	{7, 0b0011101}:  0x1a,
	{7, 0b0011100}:  0x1b,
	{7, 0b0011011}:  0x1c,
	{7, 0b0011010}:  0x1d,
	{7, 0b0011001}:  0x1e,
	{7, 0b0011000}:  0x1f,
	{7, 0b0010111}:  0x20,
	{7, 0b0010110}:  0x21,
	{7, 0b0010101}:  0x22,
	{7, 0b0010100}:  0x23,
	{7, 0b0010011}:  0x24,
	{7, 0b0010010}:  0x25,
	{7, 0b0010001}:  0x26,
	{7, 0b0010000}:  0x27,
	{7, 0b0001111}:  0x28,
	{7, 0b0001110}:  0x29,
	{7, 0b0001101}:  0x2a,
	{7, 0b0001100}:  0x2b,
	{7, 0b0001011}:  0x2c,
	{7, 0b0001010}:  0x2d,
	{7, 0b0001001}:  0x2e,
	{7, 0b0001000}:  0x2f,
	{8, 0b00001111}: 0x30,
	{8, 0b00001110}: 0x31,
	{8, 0b00001101}: 0x32,
	{8, 0b00001100}: 0x33,
	{8, 0b00001011}: 0x34,
	{8, 0b00001010}: 0x35,
	{8, 0b00001001}: 0x36,
	{8, 0b00001000}: 0x37,
	{8, 0b00000111}: 0x38,
	{8, 0b00000110}: 0x39,
	{8, 0b00000101}: 0x3a,
	{8, 0b00000100}: 0x3b,
	{8, 0b00000011}: 0x3c,
	{8, 0b00000010}: 0x3d,
	{8, 0b00000001}: 0x3e,
	{8, 0b00000000}: 0x3f,
}
