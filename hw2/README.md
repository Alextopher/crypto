# RSA

This is my implementation of the RSA cryptosystem for my cryptography class, written in go.

# Assignment

> You are reading the plaintext of our second assignment in cryptography.
For this assignment, implement your own RSA cryptosystem for text files.
Then, re-encrypt this plaintext with your own RSA keys and then submit the
encrypted text file via email along with the RSA secret keys (which are the
RSA modulus and RSA secret decryption exponent used). You may implement RSA
in any programming language of your choice. Submit the source for your RSA
program in a separate compressed file by email.

## Instructions

If you don't have the go compiler you can find download instructions at https://go.dev/doc/install. 

Running `make` is the easiest way to compile and test the program. It will encrypt and decrypt the file and run diff.

`make test` will run unit tests I wrote to verify the correctness of the program.

## Usage

```
./rsa keygen <key size> <public_key> <private_key>
./rsa encrypt <public_key> <plaintext> <ciphertext>
./rsa decrypt <private_key> <ciphertext> <plaintext>
```