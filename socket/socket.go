package main

import (
	"crypto/rand"
	"io"
	"math/big"
	"net"
)

// This is the code that handles key exchange and aes symmetric encryption
// It provides a nice interface where you just use channels to send and receive messages
// The protocol is first we send 4 bytes of message length followed by the actual message

// The wrapped socket connection
type Socket struct {
	conn net.Conn

	// send and receive channels
	send chan []byte
	recv chan []byte

	// Our ElGamal keys
	private *ElGamalPrivateKey

	// The other client's ElGamal keys
	public *ElGamalPublicKey

	// There are seperate aes ciphers for sending messages and recieving messages
	aesSend *AES
	aesRecv *AES

	// A flag that indicates if the handshake has been completed
	handshook bool
}

func NewSocket(conn net.Conn) *Socket {
	// Create a new socket
	s := &Socket{
		conn: conn,
		send: make(chan []byte),
		recv: make(chan []byte),
	}

	// force the use of 512 bit prime
	s.private, _ = Keygen(512)

	// Start the send and recieve goroutines
	go s.sendLoop()
	go s.recvLoop()

	s.Handshake()

	return s
}

// Performs a handshake to key exchange to an aes cipher
func (s *Socket) Handshake() {
	// First share ElGamal public keys
	// p, g, h
	s.send <- s.private.public.p.Bytes()
	s.send <- s.private.public.g.Bytes()
	s.send <- s.private.public.h.Bytes()

	// Read the other client's ElGamal public keys
	p := new(big.Int).SetBytes(<-s.recv)
	g := new(big.Int).SetBytes(<-s.recv)
	h := new(big.Int).SetBytes(<-s.recv)

	// Create the other client's ElGamal public key
	s.public = &ElGamalPublicKey{p, g, h}

	// Choose a random 32 bytes (16 bytes per key) to act as our half of the shared secret
	ourSecret := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, ourSecret)
	if err != nil {
		panic("could not generate random secret")
	}

	// Encrypt the shared secret with the other client's ElGamal public key
	ciphers, err := s.public.Encrypt(ourSecret)
	if err != nil {
		panic("could not encrypt ourSecret")
	}

	// First send how many ciphers we are sending
	s.send <- []byte{byte(len(ciphers))}

	// Send the ciphers [shared, cipher], [shared, cipher]...
	for _, cipher := range ciphers {
		s.send <- cipher.shared.Bytes()
		s.send <- cipher.ciphertext.Bytes()

		// Encode the size of the cipher block
		sizeBytes := make([]byte, 4)
		sizeBytes[0] = byte(cipher.size >> 24)
		sizeBytes[1] = byte(cipher.size >> 16)
		sizeBytes[2] = byte(cipher.size >> 8)
		sizeBytes[3] = byte(cipher.size)

		// Send the size of the cipher block
		s.send <- sizeBytes
	}

	// Read the number of ciphers we are receiving
	numCiphers := int((<-s.recv)[0])

	receivedCiphers := make([]*ElGamalCipherText, numCiphers)

	// Read the ciphers [shared, cipher], [shared, cipher]...
	for i := 0; i < numCiphers; i++ {
		shared := new(big.Int).SetBytes(<-s.recv)
		ciphertext := new(big.Int).SetBytes(<-s.recv)

		// Read the size of the cipher block
		sizeBytes := <-s.recv
		size := int(sizeBytes[0])<<24 | int(sizeBytes[1])<<16 | int(sizeBytes[2])<<8 | int(sizeBytes[3])

		// Create the cipher
		receivedCiphers[i] = &ElGamalCipherText{shared, ciphertext, size}
	}

	// Decrypt our friend's secret with our ElGamal private key
	friendSecret, err := s.private.Decrypt(receivedCiphers)
	if err != nil {
		panic("could not decrypt friendSecret")
	}

	// The keys are the xor of the shared secrets
	sharedSecret := make([]byte, 32)
	for i := 0; i < 32; i++ {
		sharedSecret[i] = ourSecret[i] ^ friendSecret[i]
	}

	// Check who's secret is larger to determine which key to use for sending and recieving
	ourInt := new(big.Int).SetBytes(ourSecret)
	friendInt := new(big.Int).SetBytes(friendSecret)

	if ourInt.Cmp(friendInt) > 0 {
		s.aesSend, err = NewAES(sharedSecret[:16])
		if err != nil {
			panic("could not create aesSend")
		}

		s.aesRecv, err = NewAES(sharedSecret[16:])
		if err != nil {
			panic("could not create aesRecv")
		}
	} else if ourInt.Cmp(friendInt) < 0 {
		s.aesSend, err = NewAES(sharedSecret[16:])
		if err != nil {
			panic("could not create aesSend")
		}

		s.aesRecv, err = NewAES(sharedSecret[:16])
		if err != nil {
			panic("could not create aesRecv")
		}
	} else {
		panic("our secret is the same as our friend's secret. Good luck with the lottery tonight.")
	}

	// Set the handshake complete flag
	s.handshook = true
}

func (s *Socket) sendLoop() {
	for {
		// Get the next message to send
		msg := <-s.send

		// Get the length of the message
		length := len(msg)

		// Encode the length of the message in 8 bytes
		lengthBytes := make([]byte, 8)
		lengthBytes[0] = byte(length >> 56)
		lengthBytes[1] = byte(length >> 48)
		lengthBytes[2] = byte(length >> 40)
		lengthBytes[3] = byte(length >> 32)
		lengthBytes[4] = byte(length >> 24)
		lengthBytes[5] = byte(length >> 16)
		lengthBytes[6] = byte(length >> 8)
		lengthBytes[7] = byte(length)

		// Send the length of the message
		s.conn.Write(lengthBytes)

		// If the message is not a multiple of 16 bytes, pad it with 0's
		if length%16 != 0 {
			msg = append(msg, make([]byte, 16-length%16)...)
		}

		// Behavior changes depending on if we have completed the handshake
		if s.handshook {
			cipher := make([]byte, len(msg))

			// Encrypt every 16 bytes of the message
			for i := 0; i < length; i += 16 {
				// Encrypt the message
				s.aesSend.Encrypt(cipher[i:i+16], msg[i:i+16])
			}

			// Send the encrypted message
			s.conn.Write(cipher)
		} else {
			// Send the message
			s.conn.Write(msg)
		}
	}
}

func (s *Socket) recvLoop() {
	for {
		// read 8 bytes for the length of the next message
		lengthBytes := make([]byte, 8)
		n, err := s.conn.Read(lengthBytes)
		if err != nil {
			if err == io.EOF {
				// The other side has closed the connection
				return
			}
			panic(err)
		}
		if n != 8 {
			panic("could not read 4 bytes for message length")
		}

		// Decode the length of the message
		length := int(lengthBytes[0])<<56 | int(lengthBytes[1])<<48 | int(lengthBytes[2])<<40 | int(lengthBytes[3])<<32 | int(lengthBytes[4])<<24 | int(lengthBytes[5])<<16 | int(lengthBytes[6])<<8 | int(lengthBytes[7])

		// the network length is divisible by 16
		var networkLength int
		if length%16 != 0 {
			networkLength = length + 16 - length%16
		} else {
			networkLength = length
		}

		// Read the message
		msg := make([]byte, networkLength)
		n, err = s.conn.Read(msg)
		if err != nil {
			if err == io.EOF {
				// The other side has closed the connection
				return
			}
			panic(err)
		}
		if n != networkLength {
			panic("Incomplete message")
		}

		// Behavior changes depending on if we have completed the handshake
		if s.handshook {
			// Decrypt every 16 bytes of the message
			for i := 0; i < length; i += 16 {
				// Decrypt the message
				s.aesRecv.Decrypt(msg[i:i+16], msg[i:i+16])
			}
		}

		s.recv <- msg[:length]
	}
}
