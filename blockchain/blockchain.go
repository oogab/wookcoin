package blockchain

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

var b *blockchain

func GetBlockChain() *blockchain {
	if b == nil {
		b = &blockchain{}
	}
	return b
}

// func (b *blockchain) getLastHash() string {
// 	if len(b.blocks) > 0 {
// 		return b.blocks[len(b.blocks)-1].hash
// 	}
// 	return ""
// }

// func (b *blockchain) addBlock(data string) {
// 	newBlock := block{data, "", b.getLastHash()}
// 	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash))
// 	newBlock.hash = fmt.Sprintf("%x", hash)
// 	b.blocks = append(b.blocks, newBlock)
// }

// func (b *blockchain) listBlocks() {
// 	for _, block := range b.blocks {
// 		fmt.Printf("Data: %s\n", block.data)
// 		fmt.Printf("Hash: %s\n", block.hash)
// 		fmt.Printf("Prev Hash: %s\n", block.prevHash)
// 	}
// }

// B1
// 	b1Hash = (data + "")

// B2
// 	b2Hash = (data + b1Hash)

// B3
// 	b3Hash = (data + b2Hash)
