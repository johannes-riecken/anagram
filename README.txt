How to build and run the exercise solution
==========================================

Summary: Run `go run anagrams.go sample.txt`
For parallel processing, first split the file into pieces, for example using the
POSIX split utility and then give all the pieces on the command line. Example:
`split /usr/share/dict/words && go run anagrams.go x??`

Details: In order to run the exercise solution, Go must be installed from
https://golang.org/ and the `go` tool's directory  should be added to the PATH
environment variable. The solution can also be built by running `go build
anagrams.go` and then run with `./anagrams sample.txt` or `./anagrams
sample00.txt sample01.txt ...`.
