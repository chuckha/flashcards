chapters := $(wildcard chapter*.txt)

vocab.txt: $(chapters)
	go run build.go $(chapters)
	