# Socket

## Assignment
> CS456/CS556/MA456 Cryptography Assignment 3 (Spring 2022)
>
>Submit this plaintext unencrypted by email along with all the ElGamal keys used for the encryption.
Write a simple socket-based chat program that uses ElGamal with DES/AES (your choice) in a programming 
language of your choice. Allow the user to specify the bit-length for the prime and the secret key for 
DES/AES. Implement a system that encrypts single ASCII characters. Arrange an online demonstration of 
your program. You may work in groups of three.

## What I did

I choose to implement AES-128 by reading the [spec](https://nvlpubs.nist.gov/nistpubs/fips/nist.fips.197.pdf). Near the end of the spec there are all kinds of testing values and I checked those against my implementation in [aes_test.go](aes_test.go). The socket based communication protocol is in [socket.go](socket.go). It provides a very nice implemenation where I can use channels to send messages back and forth. I do not send individual characters. 

## Usage

First build: `go build`

Then start a server: `./socket -server 0.0.0.0:8000`

And then a client: `./socket 127.0.0.1:8000`