package split

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestGenRandBin(t *testing.T) {
	filename := "/tmp/a.bin"
	var size int64 = 200
	genRandBinHelper(t, filename, size)
}

func genRandBinHelper(t *testing.T, filename string, size int64) {
	err := GenRandBin(filename, size)
	assert.Equal(t, nil, err)
	fileinfo, err := os.Stat(filename)
	assert.Equal(t, nil, err)
	assert.Equal(t, fileinfo.Size(), size)
}

func TestSplit(t *testing.T) {
	filename := "/tmp/b.bin"
	var size int64 = 201
	var partSize int64 = 100
	outDir := "/tmp"

	splitHelper(t, filename, size, partSize, outDir)
}

func splitHelper(t *testing.T, filename string, size int64, partSize int64, outDir string) []string {
	genRandBinHelper(t, filename, size)
	parts, err := Split(filename, partSize, outDir)
	assert.Equal(t, nil, err)
	fmt.Println(parts)
	return parts
}

func TestMerge(t *testing.T) {
	filename := "/tmp/c.bin"
	var size int64 = 201
	var partSize int64 = 100
	outDir := "/tmp"

	parts := splitHelper(t, filename, size, partSize, outDir)
	dstFilename := "/tmp/c_merge.bin"
	err := Merge(parts, dstFilename)
	assert.Equal(t, nil, err)

	equal, err := ContentEqual(filename, dstFilename)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, equal)
}

func Test_byteReaderDiff(t *testing.T) {
	type args struct {
		left  io.ByteReader
		right io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			args: args{
				left:  bytes.NewBufferString("abc"),
				right: bytes.NewBufferString("abc"),
			},
			want:    true,
			wantErr: false,
		},
		{
			args: args{
				left:  bytes.NewBufferString("abc"),
				right: bytes.NewBufferString("abd"),
			},
			want:    false,
			wantErr: false,
		},
		{
			args: args{
				left:  bytes.NewBufferString("ab"),
				right: bytes.NewBufferString("abd"),
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := byteReaderDiff(tt.args.left, tt.args.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("byteReaderDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("byteReaderDiff() got = %v, want %v", got, tt.want)
			}
		})
	}
}
