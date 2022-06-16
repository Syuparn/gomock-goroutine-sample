package main

import (
	"context"

	"github.com/syuparn/gomock-goroutine-sample/proto"
)

type personHandler struct {
	client proto.PersonClient
}

func NewPersonHandler(client proto.PersonClient) *personHandler {
	return &personHandler{client: client}
}

func (h *personHandler) GetName(ctx context.Context, personID int64) (string, error) {
	req := &proto.GetRequest{
		Id: personID,
	}

	res, err := h.client.Get(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetName(), nil
}
