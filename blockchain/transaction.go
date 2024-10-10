package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value int // Amount of coins

	// The public key of the recipient (in this
	// simple implementation, it's just an address)
	PubKey string
}

type TxInput struct {
	ID  []byte // The ID of the transaction containing the output we're spending
	Out int    // The index of the output in the referenced transaction

	// The signature to unlock the output (in this simple implementation,
	// it's just the sender's address)
	Sig string
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Coinbase transactions are special transactions that create new coins
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	// Coinbase input has no real input, so we use empty values
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to} // Reward is set to 100 coins

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accmulated, validOutputs := chain.FindSpendableOutputs(from, amount)

	if accmulated < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	// If there's a change, send it back to the sender
	if accmulated > amount {
		outputs = append(outputs, TxOutput{accmulated - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
