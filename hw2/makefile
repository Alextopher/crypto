all:
	go build -o rsa
	./rsa keygen 1024 pub.key priv.key
	./rsa encrypt pub.key message.txt cipher.txt
	./rsa decrypt priv.key cipher.txt plaintext.txt
	diff message.txt plaintext.txt

# unit tests for my package
test:
	go test -v ./...