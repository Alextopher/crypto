all:
	go build -o lfsr
	rm key cipher.txt plaintext.txt
	./lfsr keygen key
	./lfsr mix key message.txt cipher.txt
	./lfsr mix key cipher.txt plaintext.txt
	diff message.txt plaintext.txt