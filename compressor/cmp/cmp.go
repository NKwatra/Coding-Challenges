package cmp

import (
	"errors"
	"fmt"
	"nkwatra/compressor/utils"
	"strconv"
	"strings"
)

func countFrequencies(content string) map[rune]int {
	frequencies := make(map[rune]int)
	for _, letter := range content {
		prevFreq := frequencies[letter]
		frequencies[letter] = prevFreq + 1
	}
	return frequencies
}

func generateHuffmanTree(frequencies map[rune]int) utils.TreeNode[rune] {
	q := utils.NewMinPQueue[utils.TreeNode[rune]]()
	for letter, freq := range frequencies {
		node := utils.NewTreeNode[rune](freq)
		node.SetVal(letter)
		q.Add(*node)
	}
	for q.Size() > 1 {
		left, _ := q.Poll()
		right, _ := q.Poll()
		parent := utils.NewTreeNode[rune](left.Weight() + right.Weight())
		parent.SetLeft(&left)
		parent.SetRight(&right)
		q.Add(*parent)
	}
	root, _ := q.Poll()
	return root
}

func buildEncodingTable(root utils.TreeNode[rune], codes map[rune]string, buffer []string) {
	left := root.Left()
	right := root.Right()
	if left == nil && right == nil {
		codes[root.Val()] = strings.Join(buffer, "")
		return
	}
	if left != nil {
		buildEncodingTable(*left, codes, append(buffer, "0"))
	}

	if right != nil {
		buildEncodingTable(*right, codes, append(buffer, "1"))
	}
}

func serializeHeader(codes map[rune]string, extraBits int) string {
	var buffer strings.Builder
	buffer.WriteString("$")
	for letter, code := range codes {
		buffer.WriteString(string(letter) + "\xbd" + code + "\xcd")
	}
	if extraBits > 0 {
		buffer.WriteString(fmt.Sprint(extraBits))
	}
	headerLength := buffer.Len()
	return fmt.Sprint(headerLength) + buffer.String()
}

func compressText(original string, codes map[rune]string) (string, int) {
	compressed := make([]byte, 0)
	var lastByte byte
	byteIndex := 7
	extraBits := 0
	for _, letter := range original {
		code := codes[letter]
		for _, bit := range code {
			if bit == '1' {
				lastByte |= byte(1) << byteIndex
			}
			byteIndex--
			if byteIndex < 0 {
				compressed = append(compressed, lastByte)
				lastByte = 0
				byteIndex = 7
			}
		}
	}
	if byteIndex >= 0 {
		extraBits = byteIndex + 1
		compressed = append(compressed, lastByte)
	}
	return string(compressed), extraBits
}

func Compress(original string) string {
	frequencies := countFrequencies(original)
	codes := make(map[rune]string)
	if len(frequencies) > 1 {
		root := generateHuffmanTree(frequencies)
		buffer := make([]string, 0, 10)
		buildEncodingTable(root, codes, buffer)
	} else {
		for key := range frequencies {
			codes[key] = "0"
		}
	}
	var compressedFile strings.Builder
	compressed, extra := compressText(original, codes)
	header := serializeHeader(codes, extra)
	compressedFile.WriteString(header)
	compressedFile.WriteString(compressed)
	return compressedFile.String()
}

func Expand(compressed string) string {
	startIndex := strings.IndexRune(compressed, '$')
	headerLength, err := strconv.Atoi(compressed[:startIndex])
	if err != nil {
		panic(errors.New("compressor: compressed file is not valid"))
	}
	header := compressed[startIndex : startIndex+headerLength]
	codes, extraBits := deserializeHeader(header)
	text := compressed[startIndex+headerLength:]
	return deserializeText(text, extraBits, codes)
}

func deserializeHeader(header string) (map[string]rune, int) {
	allRunes := []rune(header)
	extraBits := allRunes[len(allRunes)-1]
	codeTable := make(map[string]rune)
	codes := allRunes[1 : len(allRunes)-1]
	for i := 0; i < len(codes); i++ {
		key := codes[i]
		i += 2
		start := i
		for i < len(codes) && codes[i] != 0xFFFD {
			i++
		}
		value := string(codes[start:i])
		codeTable[value] = key
	}
	return codeTable, int(extraBits - '0')
}

func deserializeText(text string, extraBits int, codes map[string]rune) string {
	rawData := []byte(text)
	var decoded strings.Builder
	var buffer strings.Builder
	for i := 0; i < len(rawData)-1; i++ {
		for j := 7; j >= 0; j-- {
			and := rawData[i] & (1 << j)
			if and > 0 {
				buffer.WriteRune('1')
			} else {
				buffer.WriteRune('0')
			}
			letter, exists := codes[buffer.String()]
			if exists {
				decoded.WriteRune(letter)
				buffer.Reset()
			}
		}
	}
	for i := 7; i >= extraBits; i-- {
		and := rawData[len(rawData)-1] & (1 << i)
		if and > 0 {
			buffer.WriteRune('1')
		} else {
			buffer.WriteRune('0')
		}
		letter, exists := codes[buffer.String()]
		if exists {
			decoded.WriteRune(letter)
			buffer.Reset()
		}
	}
	return decoded.String()
}
