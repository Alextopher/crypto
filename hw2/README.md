# RSA

This is my implementation of the RSA cryptosystem for my cryptography class, written in go.

**Assignment**:

> You are reading the plaintext of our second assignment in cryptography.
For this assignment, implement your own RSA cryptosystem for text files.
Then, re-encrypt this plaintext with your own RSA keys and then submit the
encrypted text file via email along with the RSA secret keys (which are the
RSA modulus and RSA secret decryption exponent used). You may implement RSA
in any programming language of your choice. Submit the source for your RSA
program in a separate compressed file by email.

## Instructions

If you don't have the go compiler you can find download instructions at https://go.dev/doc/install. 

Start by running `make`. This will compile and symlink the binaries, `keygen`, `encrypt`, `decrypt`.

To generate key pairs run the command of length bits `./keygen <bits>`, this will create two files `private.key` and `public.key`.

I haven't written much of a cli for this, everything is done using redirects of stdin and stdout. To encrypt a file with a public key you can run the following command `./encrypt public.key < plaintext.txt > cipher.txt`.

To decrypt a file with a private key you can run the following command `./decrypt private.key < cipher.txt > plaintext.txt`.

## Some technical stuff

The key file format is a text file containing the numbers of the modulus and the exponent separated by a newline.

A generated cipher text is just the numeric result after encryption encoded in ascii. (it's just a number) This means that's encrypting and decrypting are not the same operation.