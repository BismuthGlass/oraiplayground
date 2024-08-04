package services

import (
	"context"
	"crow/orai"
	"crow/oraiplayground/utils"
	"crow/oraiplayground/config"
	"errors"
	"sync"
	"time"
)

const AiServerCtxKey = utils.CtxKey("ServiceAiServer")

type AiServiceRequestClientInfo struct {
	Id int64 `json:"id"`
}

type AiServiceResponse struct {
	Prompt string `json:"prompt"`
	Result string `json:"response"`
	Err    error  `json:"err"`
}

type AiServiceClient struct {
	Id          int64
	Model       string
	Parameters  orai.Parameters
	Prompt      string
	Canceled    bool
	Timestamp   time.Time
	Output      chan AiServiceResponse
	Ctx         context.Context
	CtxCancelFn context.CancelFunc
}

type AiServer struct {
	OrCon   orai.OpenRouter
	Inbound chan *AiServiceClient
	Lock    sync.Mutex
	LastId  int64
	Pending []AiServiceClient
}

func NewAiServer() AiServer {
	return AiServer{
		OrCon: orai.New(config.ApiKey),
		Inbound: make(chan *AiServiceClient, 10),
	}
}

func (m *AiServer) Run() {
	for {
		req := <-m.Inbound
		result, err := m.OrCon.Prompt(req.Ctx, req.Model, req.Parameters, req.Prompt)
		if err != nil {
			req.Output <- AiServiceResponse{
				Err: err,
			}
		} else {
			req.Output <- AiServiceResponse{
				Prompt: req.Prompt,
				Result: result.Choices[0].Text,
			}
		}

		m.Lock.Lock()
		m.removeRequest(req.Id)
		m.Lock.Unlock()
	}
}

func (m *AiServer) findClientIndex(id int64) int {
	for i, c := range m.Pending {
		if c.Id == id {
			return i
		}
	}
	return -1
}

func (m *AiServer) findClient(id int64) *AiServiceClient {
	for _, c := range m.Pending {
		if c.Id == id {
			return &c
		}
	}
	return nil
}

func (m *AiServer) removeRequest(id int64) error {
	index := m.findClientIndex(id)
	if index < 0 {
		return errors.New("invalid request")
	}
	m.Pending = append(m.Pending[:index], m.Pending[index+1:]...)
	return nil
}

func (m *AiServer) RequestChannel(id int64) chan AiServiceResponse {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	client := m.findClient(id)
	if client == nil {
		return nil
	}
	return client.Output
}

func (m *AiServer) CancelRequest(id int64) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	client := m.findClient(id)
	if client == nil {
		return errors.New("invalid request")
	}
	client.CtxCancelFn()
	m.removeRequest(id)
	return nil
}

func (m *AiServer) IssueRequest(model string, parameters orai.Parameters, prompt string) AiServiceRequestClientInfo {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	newId := m.LastId + 1
	m.LastId++
	ctx, cancel := context.WithCancel(context.Background())
	client := AiServiceClient{
		Id: newId,
		Model: model,
		Parameters: parameters,
		Prompt: prompt,
		Output: make(chan AiServiceResponse),
		Ctx: ctx,
		CtxCancelFn: cancel,
	}
	m.Pending = append(m.Pending, client)
	m.Inbound <- &client
	return AiServiceRequestClientInfo{
		Id: newId,
	}
}
