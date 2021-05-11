package split

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
)

// base64 chars
const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func randomChar() byte {
	i := rand.Intn(len(encodeStd))
	return encodeStd[i]
}

// GenRandBin generate random binary file
func GenRandBin(filename string, size int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	var i int64
	for i = 0; i < size; i++ {
		err = w.WriteByte(randomChar())
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Split split file into n parts, each has at most x size
func Split(filename string, partSize int64, outDir string) ([]string, error) {
	srcFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return nil, err
	}

	var filelist []string
	l := int(fileInfo.Size() / partSize)
	if fileInfo.Size()%partSize != 0 {
		l++
	}
	for i := 0; i < l; i++ {
		dstFilename := filepath.Join(outDir, filepath.Base(filename)+fmt.Sprintf("_%d_part", i))
		filelist = append(filelist, dstFilename)

		dstFile, err := os.Create(dstFilename)
		if err != nil {
			return nil, err
		}

		n := partSize
		if i == l-1 {
			n = fileInfo.Size() % partSize
		}

		_, err = io.CopyN(dstFile, srcFile, n)
		if err != nil {
			return nil, err
		}

		err = dstFile.Close()
		if err != nil {
			return nil, err
		}
	}

	return filelist, nil
}

func Merge(filenameList []string, dstFilename string) error {
	sort.Strings(filenameList)
	dstFile, err := os.Create(dstFilename)
	if err != nil {
		return err
	}

	for _, filename := range filenameList {
		srcFile, err := os.Open(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
		err = srcFile.Close()
		if err != nil {
			return err
		}
	}

	return dstFile.Close()
}

func ContentEqual(leftFilename, rightFilename string) (bool, error) {
	leftFile, err := os.Open(leftFilename)
	if err != nil {
		return false, err
	}
	rightFile, err := os.Open(rightFilename)
	if err != nil {
		return false, err
	}
	defer func() {
		err := leftFile.Close()
		if err != nil {
			logrus.Error(err)
		}
		err = rightFile.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	return byteReaderDiff(bufio.NewReader(leftFile), bufio.NewReader(rightFile))
}

func byteReaderDiff(left, right io.ByteReader) (bool, error) {
	offset := 0
	for {
		aA, errA := left.ReadByte()
		bB, errB := right.ReadByte()
		if aA != bB {
			logrus.Infof("diff at offset %d, left: %v, right: %v", offset, aA, bB)
			return false, nil
		}

		// both reach to the end
		if errA == io.EOF && errB == io.EOF {
			return true, nil
		}
		if errA == io.EOF || errB == io.EOF {
			return false, nil
		}
		// other errors
		if errA != nil {
			return false, errA
		}
		if errB != nil {
			return false, errB
		}

		offset++
	}
}
