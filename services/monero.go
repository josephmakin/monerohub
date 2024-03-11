package services

import (
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

var Monero wallet.Client

func InitMonero(host, port string) error {
	address := "http://" + host + ":" + port + "/json_rpc"

	Monero = wallet.New(wallet.Config{
		Address: address,
	})

	_, err := Monero.GetBalance(&wallet.RequestGetBalance{AccountIndex: 0})
	if err != nil {
		return err
	}

	return nil
}
