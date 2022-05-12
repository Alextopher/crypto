# Ellgamal

The private `a` value was found by solving assignment 4. So in [a3.keys](a3.keys) we have all of the information to just decrypt the file. In [a3.cipher](a3.cipher) the format is 

```
cipher,halfmask
cipher,halfmask
cipher,halfmask
cipher,halfmask
...
```

Every cipher message encodes a single character.

## Resulting message

> CS456/CS556/MA456 Cryptography Assignment 3 (Spring 2022)
>
>Submit this plaintext unencrypted by email along with all the ElGamal keys used for the encryption.
Write a simple socket-based chat program that uses ElGamal with DES/AES (your choice) in a programming 
language of your choice. Allow the user to specify the bit-length for the prime and the secret key for 
DES/AES. Implement a system that encrypts single ASCII characters. Arrange an online demonstration of 
your program. You may work in groups of three.

My submission of this is in [socket](../socket/) folder.
