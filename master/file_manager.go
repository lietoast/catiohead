package master

import (
	"encoding/json"
	"sync"

	"github.com/tidwall/btree"
)

const metadataFilePath = "./metadata.json" // test environment

type FileMetadata struct {
	Mtx        *sync.RWMutex
	Filename   string                 `json:"filename"`   // file path
	Permission int16                  `json:"permission"` // access permission
	Chuncks    map[int]ChunckMetadata `json:"chuncks"`    // index -> chunck
}

type ChunckMetadata struct {
	Mtx      *sync.RWMutex
	Primary  int   `json:"primary"`  // server id of primary chunckserver
	Replicas []int `json:"replicas"` // set of replica locations
}

var fileTree btree.Map[string, FileMetadata]

func init() {
	fileJson := readJSON(metadataFilePath)
	files := make([]FileMetadata, 0)
	err := json.Unmarshal([]byte(fileJson), &files)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileTree.Set(file.Filename, file)
	}
}
