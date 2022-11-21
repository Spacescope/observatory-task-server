package chainnotifyclient

import (
	"fmt"

	"github.com/Spacescore/observatory-task/pkg/errors"
	"github.com/imroc/req/v3"
)

var (
	reqClient *req.Client
)

func init() {
	reqClient = req.C()
}

type ErrResponse struct {
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

// TopicSignIn register topic
func TopicSignIn(host string, topic string) error {
	return nil
	params := map[string]string{
		"topic": topic,
	}
	resp := reqClient.Post(fmt.Sprintf("%s/api/v1/topic", host)).SetBodyJsonMarshal(params).Do()
	if resp.Err != nil {
		return resp.Err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.IsError() {
		var errResponse ErrResponse
		if err := resp.UnmarshalJson(&errResponse); err != nil {
			return err
		}
		return errors.New(errResponse.Message)
	}
	return nil
}

// ReportTipsetState report chain notify server task state
func ReportTipsetState(host string, topic string, height, version, state, notFoundState int, desc string) error {
	return nil
	params := map[string]interface{}{
		"topic":           topic,
		"tipset":          height,
		"version":         version,
		"state":           state,
		"not_found_state": notFoundState,
		"description":     desc,
	}
	resp := reqClient.Post(fmt.Sprintf("%s/api/v1/task_state", host)).SetBodyJsonMarshal(params).Do()
	if resp.Err != nil {
		return resp.Err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.IsError() {
		var errResponse ErrResponse
		if err := resp.UnmarshalJson(&errResponse); err != nil {
			return err
		}
		return errors.New(errResponse.Message)
	}
	return nil
}
