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
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/upload2local/pkg/client"
	"github.com/xujiahua/upload2local/pkg/model"
	"github.com/xujiahua/upload2local/pkg/split"
	"os"

	"github.com/spf13/cobra"
)

var partSize int64
var hostServerURL string

const uploadUsage = "xx upload path/to/file"

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload local file to server",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetReportCaller(true)
			logrus.SetLevel(logrus.DebugLevel)
		}

		if len(args) == 0 {
			fmt.Println(uploadUsage)
			return
		}

		c := client.New(hostServerURL)

		for _, inputFilename := range args {
			partFiles, err := split.Split(inputFilename, partSize, os.TempDir())
			handleErr(err)

			for _, part := range partFiles {
				logrus.Infof("uploading file part %s ...\n", part)
				err = c.UploadFile(part)
				handleErr(err)
				logrus.Debugf("uploaded file part %s \n", part)
			}
			logrus.Debugf("finishing file uploading ... \n")
			err = c.Complete(&model.CompleteRequest{
				PartFiles: partFiles,
				Filename:  inputFilename,
			})
			handleErr(err)
			logrus.Infof("file %s uploaded\n", inputFilename)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().Int64VarP(&partSize, "size", "s", 1024*100, "split part size in bytes")
	uploadCmd.Flags().StringVarP(&hostServerURL, "host", "", "", "host server URL")
}
