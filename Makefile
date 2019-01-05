vocabs := $(wildcard vocab/chapter*.txt)
sentences := $(wildcard sentences/chapter*.txt)

all: vocab/vocab.txt sentences/sentences.txt

vocab/vocab.txt: $(vocabs)
	go run build.go vocab

sentences/sentences.txt: $(sentences)
	go run build.go sentences
