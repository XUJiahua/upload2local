package client

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/upload2local/pkg/model"
)

var bypassTunnelReminderHeader = map[string]string{
	"Bypass-Tunnel-Reminder": "1",
}

type Client struct {
	baseURL string
}

func New(baseURL string) Client {
	return Client{
		baseURL: baseURL,
	}
}

func (c Client) UploadFile(filename string) error {
	url := c.baseURL + "/"
	files, err := grequests.FileUploadFromDisk(filename)
	if err != nil {
		return err
	}

	response, err := grequests.Post(url, &grequests.RequestOptions{
		Files:   files,
		Headers: bypassTunnelReminderHeader,
	})
	if err != nil {
		return err
	}

	return responseHelper(response)
}

func (c Client) Complete(request *model.CompleteRequest) error {
	url := c.baseURL + "/complete"
	response, err := grequests.Post(url, &grequests.RequestOptions{
		JSON:    request,
		Headers: bypassTunnelReminderHeader,
	})
	if err != nil {
		return err
	}

	return responseHelper(response)
}

func responseHelper(response *grequests.Response) error {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		spew.Dump(response.Bytes())
	}

	// some errors may get from localtunnel
	// 404 means tunnel connection is lost
	// https://github.com/localtunnel/localtunnel/issues/221
	// 504 gateway timeout

	var resp model.Response
	err := response.JSON(&resp)
	if err != nil {
		return err
	}

	if resp.Code == model.CodeSuccess {
		return nil
	} else {
		return errors.New(resp.Message)
	}
}
