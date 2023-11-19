#### Compressor

- Custom implementation of encoding/decoding using the Huffman Encoding algorithm.
- Supports `-e` and `-d` flags to encoding and decoding respectively.
- Sample file `test.txt`(_3.3MB_) first decoded to `test-encoded.txt`(_1.9MB_) and then encoded back to `test-decoded.txt`(_3.3MB_) for testing

##### Test Examples

###### To Build & Install

```bash
cd compressor
go build && go install
```

###### Execute Post Installation

**Encoding**

```bash
compressor -e --out="test-encoded.txt" test.txt
```

**Decoding**

```bash
compressor -d --out="test-decoded.txt" "test-encoded.txt"
```

##### Concepts Used

1. Modules & Packages
2. Priority Queues
   - Min Heaps
   - Comparable
3. Generics
4. Unit Testing
5. Maps
6. Binary Trees
7. Bitwise operators
8. Non Unicode Bytes
9. Variable Length Encodings (**Huffman Encoding**)

##### Things I learned

- `comparable` is a predefined interface in go that adds the constraints that concrete implementations of comparable need to support `==` and `!=` operators. It does not require support for other comparison operators.

- Go uses `NewThing` / `MakeThing` type methods instead of constructors(_convention_). These methods can be used to provide custom zero values to user defined types.`NewThing` like constructors return a pointer while `MakeThing` like constructor return objects.

- struct and array are values in Go and if not passed via pointers they are passed by value i.e. copy is passed instead of address.

- Invalid Unicode code points (like `\xbd`) are processed as single rune in go and always return value `0xFFFD`, such code points can be used as separators in encoding text.
