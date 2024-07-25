package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"trading-ace/model"
)

func BlockNumberByTime(blockTime time.Time) (*model.BlockNumberResponse, *model.AppError) {
	url := fmt.Sprintf("https://coins.llama.fi/block/ethereum/%d", blockTime.Unix())
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, model.NewAppError(model.OperateFailed).Err(fmt.Errorf("get block by time new request error: %w", err))
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, model.NewAppError(model.OperateFailed).Err(fmt.Errorf("get block by time request error: %w", err))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, model.NewAppError(model.OperateFailed).Err(fmt.Errorf("get block by time read response error: %w", err))
	}

	ret := &model.BlockNumberResponse{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, model.NewAppError(model.OperateFailed).Err(fmt.Errorf("get block by time unmarshal response error: %w", err))
	}

	return ret, nil
}
