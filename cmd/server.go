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
	"github.com/NoahShen/gotunnelme/src/gotunnelme"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xujiahua/upload2local/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var inboxDirectory string
var port int
var localtunnel bool

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start server for receiving files",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetReportCaller(true)
			logrus.SetLevel(logrus.DebugLevel)
		}

		svr, err := server.New(inboxDirectory, port)
		handleErr(err)

		// start server
		go func() {
			fmt.Printf("local access: http://127.0.0.1:%d\n", port)
			handleErr(svr.Start())
		}()

		// start tunnel
		if localtunnel {
			go func() {
				// wait for server start
				time.Sleep(time.Second)

				t := gotunnelme.NewTunnel()
				url, err := t.GetUrl("")
				handleErr(err)
				fmt.Println("public access:", url)

				err = t.CreateTunnel(port)
				handleErr(err)

				t.StopTunnel()
			}()
		}

		signalChan := make(chan os.Signal, 1)
		signal.Notify(
			signalChan,
			syscall.SIGHUP,  // kill -SIGHUP XXXX
			syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
			syscall.SIGQUIT, // kill -SIGQUIT XXXX
		)
		<-signalChan
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&inboxDirectory, "inboxDirectory", "d", "./data", "inbox folder")
	serverCmd.Flags().IntVarP(&port, "port", "p", 1234, "server port")
	serverCmd.Flags().BoolVarP(&localtunnel, "localtunnel", "", true, "localtunnel enabled for public access")
}
