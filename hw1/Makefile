all:
	go build -o lsfr
	./lsfr mix key message.txt cipher.txt
	./lsfr mix key cipher.txt plaintext.txt
	diff message.txt plaintext.txt