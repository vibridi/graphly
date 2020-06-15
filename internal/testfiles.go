package internal

import (
	"encoding/json"
	"io/ioutil"

	"github.com/vibridi/graphly"
)

const (
	DirFlattened = "../internal/test/elk/flattened"
	DirCyclic    = "../internal/test/elk/cyclic"
)

func ReadTestFile(dir, name string) *graphly.Node {
	f, err := ioutil.ReadFile(dir + "/" + name)
	if err != nil {
		panic(err)
	}
	root := &graphly.Node{}
	if err := json.Unmarshal(f, root); err != nil {
		panic(err)
	}
	root.ID = name
	return root
}

func ReadTestFilesFirstn(dir string, n int) []*graphly.Node {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*graphly.Node, n)
	for i := 0; i < n; i++ {
		graphs[i] = ReadTestFile(dir, fs[i].Name())
	}
	return graphs
}

func ReadTestFilesRandn(dir string, n int) []*graphly.Node {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*graphly.Node, n)
	for i := 0; i < n; i++ {
		graphs[i] = ReadTestFile(dir, fs[RandInt(len(fs))].Name())
	}
	return graphs
}

func ReadTestFilesAll(dir string) []*graphly.Node {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*graphly.Node, len(fs))
	for i, f := range fs {
		graphs[i] = ReadTestFile(dir, f.Name())
	}
	return graphs
}
