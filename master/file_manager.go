package master

import (
	"encoding/json"
	"strings"
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

// comparFileMetadata is a comparison function that compares filenames and returns true
// when f's Filename is less than g's Filename
func comparFileMetadata(f, g FileMetadata) bool {
	return strings.Compare(f.Filename, g.Filename) < 0
}

type ChunckMetadata struct {
	Mtx      *sync.RWMutex
	Primary  int   `json:"primary"`  // server id of primary chunckserver
	Replicas []int `json:"replicas"` // set of replica locations
}

var fileTree *btree.BTreeG[FileMetadata]

func init() {
	fileJson := readJSON(metadataFilePath)
	files := make([]FileMetadata, 0)
	err := json.Unmarshal([]byte(fileJson), &files)
	if err != nil {
		panic(err)
	}

	fileTree = btree.NewBTreeG[FileMetadata](comparFileMetadata)

	for _, file := range files {
		fileTree.Set(file)
	}
}
