package main

import (
	"fmt"
	"log"
	"os"

	"reflect"

	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/heimdall/hashing"
	"github.com/it-chain/heimdall/key"
)

/*
This sample shows data to be transmitted
is signed and verified by ECDSA Key.
*/

func main() {

	keyManager, err := key.NewKeyManager("")
	errorCheck(err)

	defer os.RemoveAll("./.heimdall")

	// Generate key pair with ECDSA algorithm.
	pri, pub, err := keyManager.GenerateKey(key.ECDSA384)
	errorCheck(err)

	// Get key from memory of keyManager or from key files in key path of keyManager.
	pri, pub, err = keyManager.GetKey()
	errorCheck(err)

	// Convert key to PEM(byte) format.
	bytePriKey, err := pri.ToPEM()
	bytePubKey, err := pub.ToPEM()

	// Reconstruct key pair from bytes to key.
	err = keyManager.ByteToKey(bytePriKey, key.ECDSA384, key.PRIVATE_KEY)
	err = keyManager.ByteToKey(bytePubKey, key.ECDSA384, key.PUBLIC_KEY)
	errorCheck(err)

	// Get the reconstructed key pair.
	recPri, recPub, err := keyManager.GetKey()

	// Compare reconstructed key pair with original key pair.
	if reflect.DeepEqual(pri, recPri) && reflect.DeepEqual(pub, recPub) {
		print("reconstruct complete!\n")
	}

	sampleData := []byte("This is sample data from heimdall.")

	hashManager, err := hashing.NewHashManager()
	errorCheck(err)

	// Convert raw data to digest(hash value) by using SHA512 function.
	digest, err := hashManager.Hash(sampleData, nil, hashing.SHA512)
	errorCheck(err)

	authManager, err := auth.NewAuth()
	errorCheck(err)

	// AuthManager makes digest(hash value) to signature with private key.
	signature, err := authManager.Sign(pri, digest, nil)
	errorCheck(err)

	/* --------- After data transmitted --------- */

	// AuthManager verify that received data has any forgery during transmitting process by digest.
	// and verify that the received data is surely from the expected sender by public key.
	ok, err := authManager.Verify(pub, signature, digest, nil)
	errorCheck(err)

	fmt.Println(ok)

}

func errorCheck(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
