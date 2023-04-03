package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client interface {
	GetChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
}

type clientImpl struct {
	urlBase string
}

func (c clientImpl) GetChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	httpCient := http.Client{}

	body, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)

	httpReq, err := http.NewRequest("POST", fmt.Sprint(c.urlBase, "/chat", "/completions"), bodyReader)

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprint("Bearer ", os.Getenv("OPEN_AI_API_KEY")))

	resp, err := httpCient.Do(httpReq)

	if err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response ChatCompletionResponse

	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func NewClient() Client {
	return &clientImpl{
		urlBase: "https://api.openai.com/v1",
	}
}
