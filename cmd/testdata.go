/*
Copyright © 2021 xujiahua <littleguner@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xujiahua/upload2local/pkg/split"
)

var filename string
var size int64

// testdataCmd represents the testdata command
var testdataCmd = &cobra.Command{
	Use:   "testdata",
	Short: "generate binary data",
	Run: func(cmd *cobra.Command, args []string) {
		err := split.GenRandBin(filename, size)
		fmt.Printf("generate file: %s size: %d\n", filename, size)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("ok")
	},
}

func init() {
	rootCmd.AddCommand(testdataCmd)

	testdataCmd.Flags().StringVarP(&filename, "filepath", "f", "hello.bin", "filepath for generated file")
	testdataCmd.Flags().Int64VarP(&size, "size", "s", 1024, "size in bytes")
}
