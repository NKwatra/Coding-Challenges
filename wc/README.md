#### My Wc

- Custom implementation of the unix operator wc.
- Supports all four flags i.e. `-w`, `-c`, `-l`, `-m`
- Supports reading data from either a file or standard input.
- Contains a sample file `test.txt` used for testing.

##### Test Examples

```bash
    ./mywc -c test.txt
```

```bash
    ./mywc -l test.txt
```

```bash
    ./mywc test.txt
```

```bash
    cat test.txt | ./mywc -m
```

##### Concepts Used

1. Error Handling
   - Panic, Recovery & Error formatting
2. Defer
3. Reading command line flags and args
4. Reading files and standard input

##### Things I learned

- Reading the stats on standard input to determine the size i.e. number of bytes is not OS agnostic. Check [here](https://github.com/golang/go/issues/62392) for more details.
