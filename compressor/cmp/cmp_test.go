package cmp

import (
	"fmt"
	"nkwatra/compressor/utils"
	"strings"
	"testing"
)

func TestCountFrequencies(t *testing.T) {
	content := `This is a dummy Text, for tesTing frequencies.`
	expectedFreq := map[rune]int{
		'T': 3,
		'h': 1,
		'i': 4,
		's': 4,
		' ': 7,
		'a': 1,
		'd': 1,
		'u': 2,
		'm': 2,
		'y': 1,
		'e': 5,
		'x': 1,
		't': 2,
		',': 1,
		'f': 2,
		'o': 1,
		'r': 2,
		'n': 2,
		'g': 1,
		'q': 1,
		'c': 1,
		'.': 1,
	}
	freq := countFrequencies(content)
	for key, value := range freq {
		expectedValue, exists := expectedFreq[key]
		if !exists || expectedValue != value {
			t.Fatalf("Frequencies don't match. For %c, exptected:%d, found:%d", key, expectedValue, value)
		}
		delete(expectedFreq, key)
		delete(freq, key)
	}
	if len(expectedFreq) > 0 {
		notFoundKeys := ""
		for key := range expectedFreq {
			notFoundKeys += string(key) + ","
		}
		t.Fatalf("Expected runes %s not found", notFoundKeys)
	}
	if len(freq) > 0 {
		extraKeys := ""
		for key := range freq {
			extraKeys += string(key) + ","
		}
		t.Fatalf("Runes %s were not expected but found", extraKeys)
	}
}

func serialiseTree(root utils.TreeNode[rune], algo string, builder *strings.Builder) {
	if algo == "preorder" {
		if root.Left() == nil && root.Right() == nil {
			builder.WriteString("(" + fmt.Sprint(root.Weight()) + "," + string(root.Val()) + ")")
		} else {
			builder.WriteString("(" + fmt.Sprint(root.Weight()) + ")")
		}
		if root.Left() != nil {
			serialiseTree(*root.Left(), algo, builder)
		}
		if root.Right() != nil {
			serialiseTree(*root.Right(), algo, builder)
		}
	} else {
		if root.Left() != nil {
			serialiseTree(*root.Left(), algo, builder)
		}
		if root.Left() == nil && root.Right() == nil {
			builder.WriteString("(" + fmt.Sprint(root.Weight()) + "," + string(root.Val()) + ")")
		} else {
			builder.WriteString("(" + fmt.Sprint(root.Weight()) + ")")
		}
		if root.Right() != nil {
			serialiseTree(*root.Right(), algo, builder)
		}
	}
}

func TestGenerateHuffmanTree(t *testing.T) {
	frequencies := map[rune]int{
		'c': 32,
		'd': 42,
		'e': 120,
		'k': 7,
		'l': 42,
		'm': 24,
		'u': 37,
		'z': 2,
	}
	const inorder = "(120,e)(306)(37,u)(79)(42,d)(186)(42,l)(107)(32,c)(65)(2,z)(9)(7,k)(33)(24,m)"
	const inorder2 = "(120,e)(306)(37,u)(79)(42,l)(186)(42,d)(107)(32,c)(65)(2,z)(9)(7,k)(33)(24,m)"
	const preorder = "(306)(120,e)(186)(79)(37,u)(42,d)(107)(42,l)(65)(32,c)(33)(9)(2,z)(7,k)(24,m)"
	const preorder2 = "(306)(120,e)(186)(79)(37,u)(42,l)(107)(42,d)(65)(32,c)(33)(9)(2,z)(7,k)(24,m)"
	root := generateHuffmanTree(frequencies)
	var builder strings.Builder
	serialiseTree(root, "preorder", &builder)
	computedTree := builder.String()
	if preorder != computedTree && preorder2 != computedTree {
		t.Errorf("Huffman tree incorrect, expected %s, received %s", preorder, computedTree)
	}
	builder.Reset()
	serialiseTree(root, "inorder", &builder)
	computedTree = builder.String()
	if inorder != computedTree && inorder2 != computedTree {
		t.Errorf("Huffman tree incorrect, expected %s, received %s", inorder, computedTree)
	}
}

func TestBuildEncodingTable(t *testing.T) {
	frequencies := map[rune]int{
		'c': 32,
		'd': 42,
		'e': 120,
		'k': 7,
		'l': 42,
		'm': 24,
		'u': 37,
		'z': 2,
	}
	root := generateHuffmanTree(frequencies)
	codes := make(map[rune]string)
	buffer := make([]string, 0)
	buildEncodingTable(root, codes, buffer)
	code1 := map[rune]string{
		'c': "1110",
		'd': "101",
		'e': "0",
		'k': "111101",
		'l': "110",
		'm': "11111",
		'u': "100",
		'z': "111100",
	}
	code2 := map[rune]string{
		'c': "1110",
		'd': "110",
		'e': "0",
		'k': "111101",
		'l': "101",
		'm': "11111",
		'u': "100",
		'z': "111100",
	}
	for letter, code := range codes {
		if code != code1[letter] && code != code2[letter] {
			t.Errorf("Expected code %s for %c, received %s", code1[letter], letter, code)
		}
	}
}
