package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	img := "./img/asus-fx505ge-i7-8750h-8gb-1tb-128gb-gtx1050ti-win1-thumb-600x600.jpg"
	file, err := os.Open(img)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//read img
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	//hash img
	hash := sha256.New()
	hash.Write(bytes)
	md := hash.Sum(nil)
	asset_hash := hex.EncodeToString(md)
	fmt.Println(asset_hash)
}
