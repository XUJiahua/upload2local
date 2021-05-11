package client

import (
	"errors"
	"github.com/levigross/grequests"
	"github.com/xujiahua/upload2local/pkg/model"
)

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

	response, err := grequests.Post(url, &grequests.RequestOptions{Files: files})
	if err != nil {
		return err
	}

	return responseHelper(response)
}

func (c Client) Complete(request *model.CompleteRequest) error {
	url := c.baseURL + "/complete"
	response, err := grequests.Post(url, &grequests.RequestOptions{
		JSON: request,
	})
	if err != nil {
		return err
	}

	return responseHelper(response)
}

func responseHelper(response *grequests.Response) error {
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
