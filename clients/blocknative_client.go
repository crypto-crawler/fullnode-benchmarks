package clients

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/crypto-crawler/fullnode-benchmarks/constant"
	"github.com/crypto-crawler/fullnode-benchmarks/pojo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

// Forked from https://github.com/bonedaddy/go-blocknative/blob/main/client/client.go

// BlocknativeClient wraps gorilla websocket connections
type BlocknativeClient struct {
	apiKey        string
	system        string
	network       string
	fromWhiteList map[common.Address]bool
	toWhiteList   map[common.Address]bool
	conn          *websocket.Conn
	mtx           sync.RWMutex
}

// Create a blocknative websocket client.
//
// `system`, available values are: bitcoin, ethereum.
//
// `network`, available values are: main, ropsten, rinkeby, goerli, kovan,
// xdai, bsc-main, see https://docs.blocknative.com/mempool-explorer#supported-networks
func NewBlocknativeClient(apiKey string, system string, network string, fromWhiteList map[common.Address]bool, toWhiteList map[common.Address]bool) (*BlocknativeClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial("wss://api.blocknative.com/v0", nil)
	if err != nil {
		return nil, err
	}

	// this checks out connection to blocknative's api and makes sure that we connected properly
	var out pojo.ConnectResponse
	if err := conn.ReadJSON(&out); err != nil {
		return nil, err
	}
	if out.Status != "ok" {
		return nil, fmt.Errorf("failed to initialize websockets connection reason: %s", out.Reason)
	}

	client := &BlocknativeClient{
		apiKey:        apiKey,
		system:        system,
		network:       network,
		fromWhiteList: fromWhiteList,
		toWhiteList:   toWhiteList,
		conn:          conn,
	}
	err = client.initialize()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *BlocknativeClient) reconnect() error {
	conn, _, err := websocket.DefaultDialer.Dial("wss://api.blocknative.com/v0", nil)
	if err != nil {
		return err
	}

	// this checks out connection to blocknative's api and makes sure that we connected properly
	var out pojo.ConnectResponse
	if err := conn.ReadJSON(&out); err != nil {
		return err
	}
	if out.Status != "ok" {
		return fmt.Errorf("failed to initialize websockets connection reason: %s", out.Reason)
	}
	c.conn = conn

	err = c.initialize()
	if err != nil {
		return err
	}

	return nil
}

func (c *BlocknativeClient) ResetFromToList(fromWhiteList map[common.Address]bool, toWhiteList map[common.Address]bool) {
	c.fromWhiteList = fromWhiteList
	c.toWhiteList = toWhiteList
}

func (c *BlocknativeClient) getBaseMsg() map[string]interface{} {
	return map[string]interface{}{
		"categoryCode": "initialize",
		"eventCode":    "checkDappId",
		"timeStamp":    time.Now(),
		"dappId":       c.apiKey,
		"version":      "1",
		"blockchain": map[string]string{
			"system":  c.system,
			"network": c.network,
		},
	}
}

// Once a connection has been created with the Blocknative WebSocket service,
// an initialization message must be sent before any other messages so that
// the API key can be validated,
// see https://docs.blocknative.com/websocket#initialization
func (c *BlocknativeClient) initialize() error {
	initMsg := c.getBaseMsg()

	if err := c.writeJSON(&initMsg); err != nil {
		return err
	}

	var out pojo.ConnectResponse
	err := c.readJSON(&out)
	if err != nil {
		return err
	}
	if out.Status != "ok" {
		return fmt.Errorf("failed to initialize API connection reason:%v", out.Reason)
	}
	return nil
}

func (c *BlocknativeClient) checkResponse() error {
	var out pojo.ConnectResponse
	err := c.readJSON(&out)
	if err != nil {
		return err
	}
	if out.Status != "ok" {
		return fmt.Errorf("failed to initialize API connection reason:%v", out.Reason)
	}
	return nil
}

// Create a subscribe command to watch addLiquidity/addLiquidityETH transactions on Pancakeswa V2.
func createPancadeRouterCommand(baseMsg map[string]interface{}, router common.Address) map[string]interface{} {
	liquidityFilter := map[string]interface{}{
		"status":                  "pending",
		"to":                      router,
		"contractCall.methodName": []string{"addLiquidityETH", "addLiquidity"},
	}

	// see https://docs.blocknative.com/websocket#configurations
	baseMsg["categoryCode"] = "configs"
	baseMsg["eventCode"] = "put"
	baseMsg["config"] = map[string]interface{}{
		"scope":        router,
		"filters":      []map[string]interface{}{liquidityFilter},
		"watchAddress": true,
	}
	return baseMsg
}

// Create a subscribe command to watch multiple addresses.
func createAddressCommandDeprecated(baseMsg map[string]interface{}, whitelist map[common.Address]bool, to_or_from bool) map[string]interface{} {
	if len(whitelist) == 0 {
		return nil
	}
	addresses := make([]string, 0)
	for address := range whitelist {
		addresses = append(addresses, address.Hex())
	}

	keyStr := ""
	if to_or_from {
		keyStr = "to"
	} else {
		keyStr = "from"
	}
	toFilter := map[string]interface{}{
		"status": "pending",
		keyStr:   addresses,
	}

	// see https://docs.blocknative.com/websocket#configurations
	baseMsg["categoryCode"] = "configs"
	baseMsg["eventCode"] = "put"
	baseMsg["config"] = map[string]interface{}{
		"scope":        "global",
		"filters":      []map[string]interface{}{toFilter},
		"watchAddress": true,
	}
	return baseMsg
}

// Create a subscribe command to watch multiple addresses.
func createAddressCommand(baseMsg map[string]interface{}, whitelist map[common.Address]bool, to_or_from bool) []map[string]interface{} {
	if len(whitelist) == 0 {
		return nil
	}
	filters := make([]map[string]interface{}, 0)

	keyStr := ""
	if to_or_from {
		keyStr = "to"
	} else {
		keyStr = "from"
	}

	for address := range whitelist {
		toFilter := map[string]interface{}{
			"status": "pending",
			keyStr:   address,
		}

		// see https://docs.blocknative.com/websocket#configurations
		baseMsg["categoryCode"] = "configs"
		baseMsg["eventCode"] = "put"
		baseMsg["config"] = map[string]interface{}{
			"scope":        address,
			"filters":      []map[string]interface{}{toFilter},
			"watchAddress": true,
		}

		filters = append(filters, baseMsg)
	}

	return filters
}

// Create subscribe command messages.
func (c *BlocknativeClient) createSubscribeCommands() []map[string]interface{} {
	filters := make([]map[string]interface{}, 0)
	router := common.HexToAddress(constant.PANCAKESWAP_V2_ROUTER_ADDRESS)
	filters = append(filters, createPancadeRouterCommand(c.getBaseMsg(), router))

	delete(c.toWhiteList, router)
	if len(c.toWhiteList) > 0 {
		filters = append(filters, createAddressCommand(c.getBaseMsg(), c.toWhiteList, true)...)
	}

	if len(c.fromWhiteList) > 0 {
		filters = append(filters, createAddressCommand(c.getBaseMsg(), c.fromWhiteList, false)...)
	}

	return filters
}

func (c *BlocknativeClient) Subscribe(stopCh <-chan struct{}, outCh chan<- pojo.TxData) error {
	commands := c.createSubscribeCommands()

	for _, command := range commands {
		if err := c.writeJSON(&command); err != nil {
			return err
		}
		if err := c.checkResponse(); err != nil {
			return err
		}
	}

	go func() {
		for {
			select {
			case <-stopCh:
				c.close()
				return
			default:
				msg := &pojo.BlocknativeMsg{}
				if err := c.readJSON(msg); err != nil {
					if e, ok := err.(*websocket.CloseError); ok {
						switch e.Code {
						case websocket.CloseNormalClosure,
							websocket.CloseGoingAway,
							websocket.CloseNoStatusReceived:
							log.Printf("Web socket closed by client: %s", err)
							log.Println("Re-connecting...")
							err := c.reconnect()
							if err != nil {
								log.Fatal(err)
							} else {
								for _, command := range commands {
									if err := c.writeJSON(&command); err != nil {
										log.Fatal(err)
									}
									if err := c.checkResponse(); err != nil {
										log.Fatal(err)
									}
								}
							}
						default:
							log.Fatal("websocket read", err)
						}
					} else {
						log.Fatal(err)
					}
				}
				if msg.Status == "ok" && msg.Event.Transaction.Status == "pending" {
					outCh <- pojo.TxData(msg)
				}
			}
		}
	}()

	return nil
}

// ReadJSON is a wrapper around Conn:ReadJSON
func (c *BlocknativeClient) readJSON(out interface{}) error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.conn.ReadJSON(out)
}

// WriteJSON is a wrapper around Conn:WriteJSON
func (c *BlocknativeClient) writeJSON(msg interface{}) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.conn.WriteJSON(msg)
}

// Close is used to terminate our websocket client
func (c *BlocknativeClient) close() error {
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	c.conn.Close()
	return err
}
