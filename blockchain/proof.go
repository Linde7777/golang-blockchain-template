package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0s

const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) ConstructData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			Int64ToBytes(int64(nonce)),
			Int64ToBytes(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// Calculate finds a nonce that meet the difficulty
func (pow *ProofOfWork) Calculate() (int, []byte) {
	var hashInBytes [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.ConstructData(nonce)
		hashInBytes = sha256.Sum256(data)

		fmt.Printf("\r%x", hashInBytes)
		var hashInInt big.Int
		hashInInt.SetBytes(hashInBytes[:])

		if hashInInt.Cmp(pow.Target) != -1 {
			nonce += 1
		} else {
			break
		}

	}
	fmt.Println()

	return nonce, hashInBytes[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.ConstructData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func Int64ToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
