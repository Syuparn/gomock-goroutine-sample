package main

import (
	"context"

	"github.com/syuparn/gomock-goroutine-sample/proto"
)

type PersonID int64
type PersonName string

type personHandler struct {
	client proto.PersonClient
}

func NewPersonHandler(client proto.PersonClient) *personHandler {
	return &personHandler{client: client}
}

func (h *personHandler) GetName(ctx context.Context, personID PersonID) (PersonName, error) {
	req := &proto.GetRequest{
		Id: int64(personID),
	}

	res, err := h.client.Get(ctx, req)
	if err != nil {
		return "", err
	}

	return PersonName(res.GetName()), nil
}
