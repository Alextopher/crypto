package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
)

// Uses the 128-bit key expansion example included in the AES spec.
// Another example key is tested in TestFullCipher
func TestKeyExpansion(t *testing.T) {
	key := [16]byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c}

	expected := [11][16]byte{
		{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c},
		{0xa0, 0xfa, 0xfe, 0x17, 0x88, 0x54, 0x2c, 0xb1, 0x23, 0xa3, 0x39, 0x39, 0x2a, 0x6c, 0x76, 0x05},
		{0xf2, 0xc2, 0x95, 0xf2, 0x7a, 0x96, 0xb9, 0x43, 0x59, 0x35, 0x80, 0x7a, 0x73, 0x59, 0xf6, 0x7f},
		{0x3d, 0x80, 0x47, 0x7d, 0x47, 0x16, 0xfe, 0x3e, 0x1e, 0x23, 0x7e, 0x44, 0x6d, 0x7a, 0x88, 0x3b},
		{0xef, 0x44, 0xa5, 0x41, 0xa8, 0x52, 0x5b, 0x7f, 0xb6, 0x71, 0x25, 0x3b, 0xdb, 0x0b, 0xad, 0x00},
		{0xd4, 0xd1, 0xc6, 0xf8, 0x7c, 0x83, 0x9d, 0x87, 0xca, 0xf2, 0xb8, 0xbc, 0x11, 0xf9, 0x15, 0xbc},
		{0x6d, 0x88, 0xa3, 0x7a, 0x11, 0x0b, 0x3e, 0xfd, 0xdb, 0xf9, 0x86, 0x41, 0xca, 0x00, 0x93, 0xfd},
		{0x4e, 0x54, 0xf7, 0x0e, 0x5f, 0x5f, 0xc9, 0xf3, 0x84, 0xa6, 0x4f, 0xb2, 0x4e, 0xa6, 0xdc, 0x4f},
		{0xea, 0xd2, 0x73, 0x21, 0xb5, 0x8d, 0xba, 0xd2, 0x31, 0x2b, 0xf5, 0x60, 0x7f, 0x8d, 0x29, 0x2f},
		{0xac, 0x77, 0x66, 0xf3, 0x19, 0xfa, 0xdc, 0x21, 0x28, 0xd1, 0x29, 0x41, 0x57, 0x5c, 0x00, 0x6e},
		{0xd0, 0x14, 0xf9, 0xa8, 0xc9, 0xee, 0x25, 0x89, 0xe1, 0x3f, 0x0c, 0xc8, 0xb6, 0x63, 0x0c, 0xa6},
	}

	w := expandKey(key)

	for i := 0; i < len(expected); i++ {
		for j := 0; j < len(expected[i]); j++ {
			if w[i][j] != expected[i][j] {
				t.Errorf("w[%d][%d] = %02x, expected %02x", i, j, w[i][j], expected[i][j])
			}
		}
	}
}

func hexToInt(s string) int {
	n, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		panic(err)
	}
	return int(n)
}

// Takes as input a 32 hex-encoded string and returns a 16-byte AES state.
func stringToState(s string) (state [4][4]byte) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			i := 2*r + 8*c
			state[r][c] = byte(hexToInt(s[i : i+2]))
		}
	}

	return
}

func stringToBytes(s string) (b [16]byte) {
	if len(s) != 32 {
		panic("stringToBytes: string must be 32 characters long")
	}

	for i := 0; i < 16; i++ {
		b[i] = byte(hexToInt(s[i*2 : (i+1)*2]))
	}

	return
}

func TestHelperFunction(t *testing.T) {
	testString := "3243f6a8885a308d313198a2e0370734"
	expectedState := [4][4]byte{
		{0x32, 0x88, 0x31, 0xe0},
		{0x43, 0x5a, 0x31, 0x37},
		{0xf6, 0x30, 0x98, 0x07},
		{0xa8, 0x8d, 0xa2, 0x34},
	}
	expectedBytes := [16]byte{0x32, 0x43, 0xf6, 0xa8, 0x88, 0x5a, 0x30, 0x8d, 0x31, 0x31, 0x98, 0xa2, 0xe0, 0x37, 0x07, 0x34}

	if stringToState(testString) != expectedState {
		t.Errorf("stringToState failed: input = %v, expected %v", testString, expectedState)
	}

	if stringToBytes(testString) != expectedBytes {
		t.Errorf("stringToBytes failed: input = %v, expected %v", testString, expectedBytes)
	}
}

// Tests AES functions by walking through the test vectors from the NIST
func TestFullCipher(t *testing.T) {
	key := stringToBytes("000102030405060708090a0b0c0d0e0f")

	aes, err := NewAES(key[:])
	if err != nil {
		t.Errorf("NewAES failed: %v", err)
	}

	aes.state = stringToState("00112233445566778899aabbccddeeff")

	// Ensure that the first state is the same as the input
	if aes.state != stringToState("00112233445566778899aabbccddeeff") {
		t.Errorf("Initilize AES start state failed")
	}

	// Check key expansion
	if aes.roundKeys[0] != stringToBytes("000102030405060708090a0b0c0d0e0f") {
		t.Errorf("round 0 key is incorrect: %v", aes.roundKeys[0])
	}

	type round struct {
		start string
		s_box string
		s_row string
		m_col string
		k_sch string
	}

	var expectedRounds = [10]round{
		{
			start: "00102030405060708090a0b0c0d0e0f0",
			s_box: "63cab7040953d051cd60e0e7ba70e18c",
			s_row: "6353e08c0960e104cd70b751bacad0e7",
			m_col: "5f72641557f5bc92f7be3b291db9f91a",
			k_sch: "d6aa74fdd2af72fadaa678f1d6ab76fe",
		},
		{
			start: "89d810e8855ace682d1843d8cb128fe4",
			s_box: "a761ca9b97be8b45d8ad1a611fc97369",
			s_row: "a7be1a6997ad739bd8c9ca451f618b61",
			m_col: "ff87968431d86a51645151fa773ad009",
			k_sch: "b692cf0b643dbdf1be9bc5006830b3fe",
		},
		{
			start: "4915598f55e5d7a0daca94fa1f0a63f7",
			s_box: "3b59cb73fcd90ee05774222dc067fb68",
			s_row: "3bd92268fc74fb735767cbe0c0590e2d",
			m_col: "4c9c1e66f771f0762c3f868e534df256",
			k_sch: "b6ff744ed2c2c9bf6c590cbf0469bf41",
		},
		{
			start: "fa636a2825b339c940668a3157244d17",
			s_box: "2dfb02343f6d12dd09337ec75b36e3f0",
			s_row: "2d6d7ef03f33e334093602dd5bfb12c7",
			m_col: "6385b79ffc538df997be478e7547d691",
			k_sch: "47f7f7bc95353e03f96c32bcfd058dfd",
		},
		{
			start: "247240236966b3fa6ed2753288425b6c",
			s_box: "36400926f9336d2d9fb59d23c42c3950",
			s_row: "36339d50f9b539269f2c092dc4406d23",
			m_col: "f4bcd45432e554d075f1d6c51dd03b3c",
			k_sch: "3caaa3e8a99f9deb50f3af57adf622aa",
		},
		{
			start: "c81677bc9b7ac93b25027992b0261996",
			s_box: "e847f56514dadde23f77b64fe7f7d490",
			s_row: "e8dab6901477d4653ff7f5e2e747dd4f",
			m_col: "9816ee7400f87f556b2c049c8e5ad036",
			k_sch: "5e390f7df7a69296a7553dc10aa31f6b",
		},
		{
			start: "c62fe109f75eedc3cc79395d84f9cf5d",
			s_box: "b415f8016858552e4bb6124c5f998a4c",
			s_row: "b458124c68b68a014b99f82e5f15554c",
			m_col: "c57e1c159a9bd286f05f4be098c63439",
			k_sch: "14f9701ae35fe28c440adf4d4ea9c026",
		},
		{
			start: "d1876c0f79c4300ab45594add66ff41f",
			s_box: "3e175076b61c04678dfc2295f6a8bfc0",
			s_row: "3e1c22c0b6fcbf768da85067f6170495",
			m_col: "baa03de7a1f9b56ed5512cba5f414d23",
			k_sch: "47438735a41c65b9e016baf4aebf7ad2",
		},
		{
			start: "fde3bad205e5d0d73547964ef1fe37f1",
			s_box: "5411f4b56bd9700e96a0902fa1bb9aa1",
			s_row: "54d990a16ba09ab596bbf40ea111702f",
			m_col: "e9f74eec023020f61bf2ccf2353c21c7",
			k_sch: "549932d1f08557681093ed9cbe2c974e",
		},
	}

	// now there are 9 rounds of normal encryption
	for i := 0; i < 9; i++ {
		// Add round key
		aes.addRoundKey(i)
		if aes.state != stringToState(expectedRounds[i].start) {
			t.Errorf("round %d start is incorrect: %v", i+1, aes.state)
			return
		}

		// SubBytes
		aes.subBytes()
		if aes.state != stringToState(expectedRounds[i].s_box) {
			t.Errorf("round %d s_box is incorrect: %v", i+1, aes.state)
			return
		}

		// ShiftRows
		aes.shiftRows()
		if aes.state != stringToState(expectedRounds[i].s_row) {
			t.Errorf("round %d s_row is incorrect: %v", i+1, aes.state)
			return
		}

		// MixColumns
		aes.mixColumns()
		if aes.state != stringToState(expectedRounds[i].m_col) {
			t.Errorf("round %d m_col is incorrect: %v", i+1, aes.state)
			return
		}

		// Check the key schedule
		if aes.roundKeys[i+1] != stringToBytes(expectedRounds[i].k_sch) {
			t.Errorf("round %d key schedule is incorrect: %v", i+1, aes.roundKeys[i+1])
			return
		}
	}

	// The last round is special

	// Add round key
	aes.addRoundKey(9)
	if aes.state != stringToState("bd6e7c3df2b5779e0b61216e8b10b689") {
		t.Errorf("round 10 start is incorrect: %v", aes.state)
		return
	}

	// SubBytes
	aes.subBytes()
	if aes.state != stringToState("7a9f102789d5f50b2beffd9f3dca4ea7") {
		t.Errorf("round 10 s_box is incorrect: %v", aes.state)
		return
	}

	// ShiftRows
	aes.shiftRows()
	if aes.state != stringToState("7ad5fda789ef4e272bca100b3d9ff59f") {
		t.Errorf("round 10 s_row is incorrect: %v", aes.state)
		return
	}

	// Check the key schedule
	if aes.roundKeys[10] != stringToBytes("13111d7fe3944a17f307a78b4d2b30c5") {
		t.Errorf("round 10 key schedule is incorrect: %v", aes.roundKeys[10])
		return
	}

	// Final AddRoundKey
	aes.addRoundKey(10)
	if aes.state != stringToState("69c4e0d86a7b0430d8cdb78070b4c55a") {
		t.Errorf("round 10 end is incorrect: %v", aes.state)
		return
	}
}

func TestFullInverseCipher(t *testing.T) {
	key := stringToBytes("000102030405060708090a0b0c0d0e0f")

	aes, err := NewAES(key[:])
	if err != nil {
		t.Errorf("NewAES failed: %v", err)
	}

	aes.state = stringToState("69c4e0d86a7b0430d8cdb78070b4c55a")

	// Ensure that the first state is the same as the input
	if aes.state != stringToState("69c4e0d86a7b0430d8cdb78070b4c55a") {
		t.Errorf("Initilize AES start state failed")
	}

	type round struct {
		istart string
		is_row string
		is_box string
		ik_sch string
		ik_add string
	}

	var expectedRounds = [10]round{
		{
			istart: "7ad5fda789ef4e272bca100b3d9ff59f",
			is_row: "7a9f102789d5f50b2beffd9f3dca4ea7",
			is_box: "bd6e7c3df2b5779e0b61216e8b10b689",
			ik_sch: "549932d1f08557681093ed9cbe2c974e",
			ik_add: "e9f74eec023020f61bf2ccf2353c21c7",
		},
		{
			istart: "54d990a16ba09ab596bbf40ea111702f",
			is_row: "5411f4b56bd9700e96a0902fa1bb9aa1",
			is_box: "fde3bad205e5d0d73547964ef1fe37f1",
			ik_sch: "47438735a41c65b9e016baf4aebf7ad2",
			ik_add: "baa03de7a1f9b56ed5512cba5f414d23",
		},
		{
			istart: "3e1c22c0b6fcbf768da85067f6170495",
			is_row: "3e175076b61c04678dfc2295f6a8bfc0",
			is_box: "d1876c0f79c4300ab45594add66ff41f",
			ik_sch: "14f9701ae35fe28c440adf4d4ea9c026",
			ik_add: "c57e1c159a9bd286f05f4be098c63439",
		},
		{
			istart: "b458124c68b68a014b99f82e5f15554c",
			is_row: "b415f8016858552e4bb6124c5f998a4c",
			is_box: "c62fe109f75eedc3cc79395d84f9cf5d",
			ik_sch: "5e390f7df7a69296a7553dc10aa31f6b",
			ik_add: "9816ee7400f87f556b2c049c8e5ad036",
		},
		{
			istart: "e8dab6901477d4653ff7f5e2e747dd4f",
			is_row: "e847f56514dadde23f77b64fe7f7d490",
			is_box: "c81677bc9b7ac93b25027992b0261996",
			ik_sch: "3caaa3e8a99f9deb50f3af57adf622aa",
			ik_add: "f4bcd45432e554d075f1d6c51dd03b3c",
		},
		{
			istart: "36339d50f9b539269f2c092dc4406d23",
			is_row: "36400926f9336d2d9fb59d23c42c3950",
			is_box: "247240236966b3fa6ed2753288425b6c",
			ik_sch: "47f7f7bc95353e03f96c32bcfd058dfd",
			ik_add: "6385b79ffc538df997be478e7547d691",
		},
		{
			istart: "2d6d7ef03f33e334093602dd5bfb12c7",
			is_row: "2dfb02343f6d12dd09337ec75b36e3f0",
			is_box: "fa636a2825b339c940668a3157244d17",
			ik_sch: "b6ff744ed2c2c9bf6c590cbf0469bf41",
			ik_add: "4c9c1e66f771f0762c3f868e534df256",
		},
		{
			istart: "3bd92268fc74fb735767cbe0c0590e2d",
			is_row: "3b59cb73fcd90ee05774222dc067fb68",
			is_box: "4915598f55e5d7a0daca94fa1f0a63f7",
			ik_sch: "b692cf0b643dbdf1be9bc5006830b3fe",
			ik_add: "ff87968431d86a51645151fa773ad009",
		},
		{
			istart: "a7be1a6997ad739bd8c9ca451f618b61",
			is_row: "a761ca9b97be8b45d8ad1a611fc97369",
			is_box: "89d810e8855ace682d1843d8cb128fe4",
			ik_sch: "d6aa74fdd2af72fadaa678f1d6ab76fe",
			ik_add: "5f72641557f5bc92f7be3b291db9f91a",
		},
	}

	// Check key expansion
	if aes.roundKeys[10] != stringToBytes("13111d7fe3944a17f307a78b4d2b30c5") {
		t.Errorf("round 0 key is incorrect: %v", aes.roundKeys[10])
	}

	// Add round key
	aes.addRoundKey(10)

	// There are 10 rounds
	for i := 0; i < 9; i++ {
		// check the start value
		if aes.state != stringToState(expectedRounds[i].istart) {
			t.Errorf("round %d start value is incorrect: %v", i, aes.state)
		}

		// inverse mix rows
		aes.invShiftRows()
		if aes.state != stringToState(expectedRounds[i].is_row) {
			t.Errorf("round %d inverse shift rows is incorrect: %v", i, aes.state)
		}

		// inverse sub bytes
		aes.invSubBytes()
		if aes.state != stringToState(expectedRounds[i].is_box) {
			t.Errorf("round %d inverse sub bytes is incorrect: %v", i, aes.state)
		}

		// check key schedule
		if aes.roundKeys[9-i] != stringToBytes(expectedRounds[i].ik_sch) {
			t.Errorf("round %d key schedule is incorrect: %v", i, aes.roundKeys[9-i])
		}

		// add round key
		aes.addRoundKey(9 - i)
		if aes.state != stringToState(expectedRounds[i].ik_add) {
			t.Errorf("round %d add round key is incorrect: %v", i, aes.state)
		}

		// inverse mix columns
		aes.invMixColumns()
	}

	// check the start value
	if aes.state != stringToState("6353e08c0960e104cd70b751bacad0e7") {
		t.Errorf("round 10 start value is incorrect: %v", aes.state)
	}

	// inverse mix rows
	aes.invShiftRows()
	if aes.state != stringToState("63cab7040953d051cd60e0e7ba70e18c") {
		t.Errorf("round 10 inverse shift rows is incorrect: %v", aes.state)
	}

	// inverse sub bytes
	aes.invSubBytes()
	if aes.state != stringToState("00102030405060708090a0b0c0d0e0f0") {
		t.Errorf("round 10 inverse sub bytes is incorrect: %v", aes.state)
	}

	// check key schedule
	if aes.roundKeys[0] != stringToBytes("000102030405060708090a0b0c0d0e0f") {
		t.Errorf("round 10 key schedule is incorrect: %v", aes.roundKeys[0])
	}

	// add round key
	aes.addRoundKey(0)
	if aes.state != stringToState("00112233445566778899aabbccddeeff") {
		t.Errorf("round 10 add round key is incorrect: %v", aes.state)
	}
}

// Tests the same thing as TestFullCipher but without looking at intermediate values
func TestBlockEncrypt(t *testing.T) {
	key := stringToBytes("000102030405060708090a0b0c0d0e0f")

	aes, err := NewAES(key[:])
	if err != nil {
		t.Errorf("NewAES failed: %v", err)
	}

	out := aes.BlockEncrypt(stringToBytes("00112233445566778899aabbccddeeff"))

	if out != stringToBytes("69c4e0d86a7b0430d8cdb78070b4c55a") {
		t.Errorf("BlockEncrypt failed: expected: %v, found: %v", out, stringToBytes("69c4e0d86a7b0430d8cdb78070b4c55a"))
	}
}

func TestBlockDecrypt(t *testing.T) {
	key := stringToBytes("000102030405060708090a0b0c0d0e0f")

	aes, err := NewAES(key[:])
	if err != nil {
		t.Errorf("NewAES failed: %v", err)
	}

	out := aes.BlockDecrypt(stringToBytes("69c4e0d86a7b0430d8cdb78070b4c55a"))

	if out != stringToBytes("00112233445566778899aabbccddeeff") {
		t.Errorf("BlockDecrypt failed: expected: %v, found: %v", out, stringToBytes("00112233445566778899aabbccddeeff"))
	}
}

func TestE2E(t *testing.T) {
	var msg, key [16]byte
	for i := 1; i < 1000; i++ {
		rand.Read(msg[:])
		rand.Read(key[:])

		aes, err := NewAES(key[:])
		if err != nil {
			t.Errorf("NewAES failed: %v", err)
		}

		cipher := aes.BlockEncrypt(msg)
		decipher := aes.BlockDecrypt(cipher)

		if !bytes.Equal(msg[:], decipher[:]) {
			t.Errorf("E2E failed: expected: %v, found: %v", msg, decipher)
		}
	}
}
