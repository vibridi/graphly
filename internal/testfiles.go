package internal

import (
	"io/ioutil"
)

type GraphData struct {
	Data []byte
	Name string
}

func ReadTestFile(dir, name string) *GraphData {
	f, err := ioutil.ReadFile(dir + "/" + name)
	if err != nil {
		panic(err)
	}
	return &GraphData{
		Data: f,
		Name: name,
	}
}

func ReadTestFilesFirstn(dir string, n int) []*GraphData {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*GraphData, n)
	for i := 0; i < n; i++ {
		graphs[i] = ReadTestFile(dir, fs[i].Name())
	}
	return graphs
}

func ReadTestFilesRandn(dir string, n int) []*GraphData {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*GraphData, n)
	for i := 0; i < n; i++ {
		graphs[i] = ReadTestFile(dir, fs[RandInt(len(fs))].Name())
	}
	return graphs
}

func ReadTestFilesAll(dir string) []*GraphData {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*GraphData, len(fs))
	for i, f := range fs {
		graphs[i] = ReadTestFile(dir, f.Name())
	}
	return graphs
}
