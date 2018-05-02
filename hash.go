package main

import "strconv"
import "os"
import "fmt"
import "bufio"
import "encoding/hex"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	hashSize := 32 // in bytes. output is hex so twice this value will be the amount of output characters
	var err error

	if len(os.Args) < 2 {
		panic("Please provide a filename")
	}

	if len(os.Args) >= 3 {
		hashSize, err = strconv.Atoi(os.Args[2])
		check(err)
	}

	fileName := os.Args[1]
	//fmt.Println("Filename is", fileName, "hash size =", hashSize)

	f, err := os.Open(fileName)
	check(err)

	hash := hash(f, hashSize)
	fmt.Println(hex.EncodeToString(hash))
	//fmt.Println(hash)
}

func hash(filePointer *os.File, hashSize int) []byte {
	defer filePointer.Close()
	var previousByte byte = 0
	var hash []byte = make([]byte, hashSize)
	i, next, filled, offset := 0, 0, 0, 0

	for filled < hashSize {
		offset = filled
		filePointer.Seek(0, 0)
		scanner := bufio.NewScanner(filePointer)
		scanner.Split(bufio.ScanBytes)
		i = (i + offset) % hashSize
		for scanner.Scan() {
			b := scanner.Bytes()[0]
			x := next + int(previousByte) + int(b)

			hash[i] = byte(int(hash[i]) + x - (x % 8))
			next = (int(previousByte) + int(b)) / 8
			previousByte = hash[i]
			i = (i + 1) % hashSize
			filled++
		}
		if filled == 0 {
			panic("empty file detected")
		}
	}

	return hash
}
