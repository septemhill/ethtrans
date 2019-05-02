package main

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type EtherRPCClient struct {
	cli *rpc.Client
}

func (e *EtherRPCClient) Request(method string, result interface{}, args ...interface{}) error {
	return e.cli.Call(result, method, args...)
}

func (e *EtherRPCClient) Close() {
	e.cli.Close()
}

func NewEtherRPCClient(url string) *EtherRPCClient {
	e := &EtherRPCClient{}
	e.cli, _ = rpc.Dial(url)

	return e
}
