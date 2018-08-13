package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"io/ioutil"
)

func TestDebArchiver(t *testing.T) {
	za := &DebArchvier{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.deb", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "data.tar.xz")
	assert.Equal(t, ad.ModTime, int64(1485714631))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(42284))
	b, err := ioutil.ReadAll(ad.ArchiveReader)
	if err != nil {
		t.Fatal(err)
	}
}
