package client

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ybbus/jsonrpc"
)

type Client struct {
	client jsonrpc.RPCClient
}

// NewClient creates a client that uses the given RPC client.
func NewClient(host string) *Client {
	client := jsonrpc.NewClient(host)
	return &Client{client: client}
}

type CallMsg struct {
	From            QkcAddress  // the sender of the 'transaction'
	To              *QkcAddress // the destination contract (nil for contract creation)
	Gas             uint64      // if 0, the call executes with near-infinite gas
	GasPrice        *big.Int    // wei <-> gas exchange ratio
	Value           *big.Int    // amount of wei sent along with the call
	Data            []byte      // input data, usually an ABI-encoded contract method invocation
	GasTokenId      uint64
	TransferTokenId uint64
}

func toCallArg(msg *CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From.ToHex(),
		"to":   msg.To.ToHex(),
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Encode(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = hexutil.EncodeBig(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.EncodeUint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = hexutil.EncodeBig(msg.GasPrice)
	}
	if msg.GasTokenId != 0 {
		arg["gasTokenId"] = hexutil.EncodeUint64(msg.GasTokenId)
	}
	if msg.TransferTokenId != 0 {
		arg["transferTokenId"] = hexutil.EncodeUint64(msg.TransferTokenId)
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

func (c *Client) GetMinorBlockByHeight(fullShardId uint32, number *big.Int) (result *jsonrpc.RPCResponse, err error) {
	nn := new(string)
	if number == nil {
		nn = nil
	} else {
		tt := hexutil.EncodeUint64(number.Uint64())
		nn = &tt
	}
	resp, err := c.client.Call("getMinorBlockByHeight", hexutil.EncodeUint64(uint64(fullShardId)), nn, false)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (c *Client) GetMinorBlockById(blockId string) (result *jsonrpc.RPCResponse, err error) {
	resp, err := c.client.Call("getMinorBlockById", blockId, true)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (c *Client) GetRootBlockHeight() (uint64, error) {
	resp, err := c.client.Call("getRootBlockByHeight")
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}
	heightStr := resp.Result.(map[string]interface{})["height"].(string)
	height := new(big.Int).SetBytes(common.FromHex(heightStr)).Uint64()
	return height, nil
}

func (c *Client) GetRootBlockByHeight(number *big.Int) (result *jsonrpc.RPCResponse, err error) {
	resp, err := c.client.Call("getRootBlockByHeight", toBlockNumArg(number))
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (c *Client) GetTransactionById(txid *TransactionId) (result *jsonrpc.RPCResponse, err error) {
	resp, err := c.client.Call("getTransactionById", []string{txid.Hex()})
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (c *Client) GetTransactionReceipt(transactionId *TransactionId) (result *jsonrpc.RPCResponse, err error) {
	resp, err := c.client.Call("getTransactionReceipt", []string{transactionId.Hex()})
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (c *Client) GetBalance(qkcAddr *QkcAddress, token string) (balance *big.Int, err error) {
	resp, err := c.client.Call("getBalances", []string{qkcAddr.ToHex()})
	if err != nil {
		return
	}
	if resp.Error != nil {
		fmt.Println("getBalances error: ", resp.Error.Error())
		return
	}
	balances := resp.Result.(map[string]interface{})["balances"]
	for _, m := range balances.([]interface{}) {
		bInfo := m.(map[string]interface{})
		if strings.ToLower((bInfo["tokenStr"]).(string)) == token {
			return hexutil.DecodeBig(bInfo["balance"].(string))
		}
	}
	return new(big.Int).SetUint64(0), nil
}

func (c *Client) GetCode(tokenAddr *QkcAddress, number *big.Int) (string, error) {
	resp, err := c.client.Call("getCode", tokenAddr.ToHex(), toBlockNumArg(number))
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}
	code, err := hexutil.Decode(resp.Result.(string))
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *Client) SendTransaction(tx *EvmTransaction) ([]byte, error) {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Call("sendRawTransaction", common.ToHex(data))
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return common.FromHex(resp.Result.(string)), nil
}

func (c *Client) GasPrice(fullShardId uint32, tokenId uint64) (*big.Int, error) {
	resp, err := c.client.Call("gasPrice", hexutil.EncodeUint64(uint64(fullShardId)), hexutil.EncodeUint64(tokenId))
	if err != nil {
		return nil, err
	}
	price, err := hexutil.DecodeBig(resp.Result.(map[string]interface{})["result"].(string))
	if err != nil {
		return nil, err
	}
	return price, nil
}

func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}
	return arg, nil
}

func (c *Client) GetLog(fullShardId uint32, query ethereum.FilterQuery) (*jsonrpc.RPCResponse, error) {
	arg, err := toFilterArg(query)
	if err != nil {
		return nil, err
	}
	// fmt.Println("ERRRR",arg,hexutil.EncodeUint64(uint64(fullShardId)))
	result, err := c.client.Call("eth_getLogs", arg, hexutil.EncodeUint64(uint64(fullShardId)))
	// fmt.Println("err",result)
	return result, err
}

func (c *Client) GetAccountData(qkcaddr *QkcAddress, number *big.Int, includeShards bool) (map[string]interface{}, error) {
	resp, err := c.client.Call("getAccountData", qkcaddr.ToHex(), nil, includeShards)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	fullShardId := GetFullShardIdByFullShardKey(qkcaddr.FullShardKey)
	shards := resp.Result.(map[string]interface{})["shards"]
	for _, val := range shards.([]interface{}) {
		shrd := val.(map[string]interface{})
		id, err := hexutil.DecodeUint64(shrd["fullShardId"].(string))
		if err != nil {
			return nil, err
		}
		if id == uint64(fullShardId) {
			return shrd, nil
		}
	}
	return nil, errors.New("has no such account")
}

func (c *Client) networkInfo() (result *jsonrpc.RPCResponse, err error) {
	resp, err := c.client.Call("networkInfo")
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (c *Client) Call(params *CallMsg) (data []byte, err error) {
	resp, err := c.client.Call("call", toCallArg(params), nil)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return hexutil.Decode(resp.Result.(string))
}

func (c *Client) GetNonce(qkcAddr *QkcAddress) (nonce uint64, err error) {
	shrd, err := c.GetAccountData(qkcAddr, nil, true)
	if err != nil {
		return 0, err
	}
	return hexutil.DecodeUint64(shrd["transactionCount"].(string))
}

func (c *Client) NetworkID() (uint32, error) {
	resp, err := c.networkInfo()
	if err != nil {
		return 0, err
	}
	networkId, err := hexutil.DecodeUint64(resp.Result.(map[string]interface{})["networkId"].(string))
	if err != nil {
		return 0, err
	}
	return uint32(networkId), nil
}

func (c *Client) CreateTransaction(networkID uint32, nonce uint64, fromFullShardKey uint32, qkcToAddr *QkcAddress, amount *big.Int, gasLimit uint64, gasPrice *big.Int, tokenID uint64, data []byte) (tx *EvmTransaction) {
	to := new(common.Address)
	if qkcToAddr == nil {
		to = nil
	} else {
		to = &qkcToAddr.Recipient
	}

	tx = NewEvmTransaction(nonce, to, amount, gasLimit, gasPrice, fromFullShardKey, fromFullShardKey,
		TokenIDEncode("QKC"), tokenID, networkID, 0, data)
	return tx
}

func (c *Client) GetFullShardIds() ([]uint32, error) {
	resp, err := c.client.Call("getFullShardIds")
	if err != nil {
		return []uint32{}, err
	}
	data := resp.Result.([]interface{})
	res := make([]uint32, 0, len(data))
	for _, id := range data {
		fullShardId, err := hexutil.DecodeUint64(id.(string))
		if err != nil {
			return nil, err
		}
		res = append(res, uint32(fullShardId))
	}
	return res, nil
}
