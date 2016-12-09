# go-packer
Fast data encoding/decoding library that works without codegeneration on golang.

It's like gob but faster.

```
BenchmarkPackerIntEncode-8    	10000000	       230 ns/op	  38.99 MB/s	       8 B/op	       1 allocs/op
BenchmarkGobInt-8             	 2000000	       867 ns/op	  13.83 MB/s	       8 B/op	       1 allocs/op
```
