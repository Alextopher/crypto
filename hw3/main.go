package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"
)

func main() {
	p, _, _, a := readKeys()

	// Read each [halfmask,cipher] from the cipher file
	cipher, err := os.Open("a3.cipher")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read the file line by line
	scanner := bufio.NewScanner(cipher)
	for scanner.Scan() {
		// Split the line at the comma
		line := strings.Split(scanner.Text(), ",")

		// Convert the halfmask to a big.Int
		halfmask := fatalStringToBigInt(line[0])

		// Convert the cipher to a big.Int
		cipher := fatalStringToBigInt(line[1])

		// Fullmask = halfmask^a mod p
		fullmask := new(big.Int).Exp(halfmask, a, p)

		// modular inverse of fullmask
		modinv := new(big.Int).ModInverse(fullmask, p)

		// Decrypt the cipher
		plain := new(big.Int).Mul(cipher, modinv)
		plain = new(big.Int).Mod(plain, p)

		// Print the plaintext as a string
		// fmt.Printf("%v\n", plain.Bytes())
		fmt.Print(string(plain.Bytes()))
	}

}

func fatalStringToBigInt(s string) *big.Int {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		fmt.Println("Could not convert string to big.Int")
		os.Exit(1)
	}
	return i
}

func readKeys() (p, g, b, a *big.Int) {
	// string to read the file into
	allKeys, err := ioutil.ReadFile("a3.keys")
	if err != nil {
		log.Fatal(err)
	}

	// Read these values with scanf
	// ````
	// p = 210721876453307145525244242041522249222704625287866585233742799076848568768836672842060166093882067569156988516568063989215149880168912278302617373054515229746447984143923759907703287852804176684489787848475892933977081903524027048664329776063475272665499609101143618674718047781127341242141906405500532997987
	// g = 95805329712244551097483332014742094876308572033539201678436979100249984241417782239221908567858200319671801015451929958159179280796193346606666274112827687714541332805153273569580116244814211930200826941979982355268422431766290700854855830849616160823820098282877902861679758560897579269065846480445940064886
	// b = 163865859017002499555661309238469696025944892886780249095613950472397889260824251591963267593348911890320183885312693071755558405005251311028940333304598675992448729251171059610996069045270901840838992983272272818400586652845890525520509188956779458733789944330079346006118667820457237065635747224427587470489
	// a = 31969280088361358982891392020418160677752584268704094389910371947753244946933347935878625023919402452740316172441746808587155134252287220284210084660983594234250923217292376810981985895141552247204162868475624471518794404318233925012571130473382462504642943159579216921221327811371889912608373396443367093840
	// ````

	var _p, _g, _b, _a string

	// Parse the whole file with regex
	re := regexp.MustCompile(`p = (\d+)\ng = (\d+)\nb = (\d+)\na = (\d+)`)
	for _, match := range re.FindAllStringSubmatch(string(allKeys), -1) {
		_p = match[1]
		_g = match[2]
		_b = match[3]
		_a = match[4]
	}

	// Convert to big.Int
	p = fatalStringToBigInt(_p)
	g = fatalStringToBigInt(_g)
	b = fatalStringToBigInt(_b)
	a = fatalStringToBigInt(_a)

	return
}
