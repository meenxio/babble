package app

import (
	"time"

	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/babble/src/node/state"
	"github.com/mosaicnetworks/babble/src/proxy"
	"github.com/sirupsen/logrus"
)

// SocketAppProxy is the Babble side of the AppProxy interface implemented over
// a TCP/RPC connection
type SocketAppProxy struct {
	clientAddress string
	bindAddress   string

	client *SocketAppProxyClient
	server *SocketAppProxyServer

	logger *logrus.Entry
}

// NewSocketAppProxy creates a new SocketAppProxy
func NewSocketAppProxy(clientAddr string, bindAddr string, timeout time.Duration, logger *logrus.Entry) (*SocketAppProxy, error) {
	if logger == nil {
		log := logrus.New()
		log.Level = logrus.DebugLevel
		logger = logrus.NewEntry(log)
	}

	client := NewSocketAppProxyClient(clientAddr, timeout, logger)

	server, err := NewSocketAppProxyServer(bindAddr, logger)

	if err != nil {
		return nil, err
	}

	proxy := &SocketAppProxy{
		clientAddress: clientAddr,
		bindAddress:   bindAddr,
		client:        client,
		server:        server,
		logger:        logger,
	}

	go proxy.server.listen()

	return proxy, nil
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//Implement AppProxy Interface

// SubmitCh implements the AppProxy interface
func (p *SocketAppProxy) SubmitCh() chan []byte {
	return p.server.submitCh
}

// CommitBlock implements the AppProxy interface
func (p *SocketAppProxy) CommitBlock(block hashgraph.Block) (proxy.CommitResponse, error) {
	return p.client.CommitBlock(block)
}

// GetSnapshot implements the AppProxy interface
func (p *SocketAppProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return p.client.GetSnapshot(blockIndex)
}

// Restore implements the AppProxy interface
func (p *SocketAppProxy) Restore(snapshot []byte) error {
	return p.client.Restore(snapshot)
}

// OnStateChanged implements the AppProxy interface
func (p *SocketAppProxy) OnStateChanged(state state.State) error {
	return p.client.OnStateChanged(state)
}
