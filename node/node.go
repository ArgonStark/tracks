package node

import (
	"context"
	"github.com/airchains-network/decentralized-sequencer/pods"
	//"fmt"
	logs "github.com/airchains-network/decentralized-sequencer/log"
	//"strconv"
	//"strings"
	"sync"

	"github.com/airchains-network/decentralized-sequencer/blocksync"
	"github.com/airchains-network/decentralized-sequencer/node/shared"
	"github.com/airchains-network/decentralized-sequencer/p2p"
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/syndtr/goleveldb/leveldb"
)

func Start() {
	var wg1 sync.WaitGroup
	wg1.Add(2)

	go configureP2P(&wg1)
	go initializeDBAndStartIndexing(&wg1)

	wg1.Wait()
}

func configureP2P(wg *sync.WaitGroup) {
	defer wg.Done()
	p2p.P2PConfiguration()
}

func initializeDBAndStartIndexing(wg *sync.WaitGroup) {
	defer wg.Done()

	staticDB := shared.Node.NodeConnections.GetBlockDatabaseConnection()

	shared.CheckAndInitializeDBCounters(staticDB)

	latestBlock := shared.GetLatestBlock(shared.Node.NodeConnections.BlockDatabaseConnection)
	client, err := ethclient.Dial("http://192.168.1.106:8545")
	if err != nil {
		logs.Log.Error("Error in connecting to the network")
		return
	}

	var ctx context.Context
	ctx = context.Background()
	var wgnm *sync.WaitGroup
	wgnm = &sync.WaitGroup{}
	wgnm.Add(2)

	latestBatch := shared.GetLatestBatchIndex(staticDB)

	//go configureP2P(wgnm)
	go blocksync.StartIndexer(wgnm, client, ctx, shared.Node.NodeConnections.BlockDatabaseConnection, shared.Node.NodeConnections.TxnDatabaseConnection, latestBlock)
	go pods.BatchGeneration(wgnm, client, ctx, shared.Node.NodeConnections.StaticDatabaseConnection, shared.Node.NodeConnections.TxnDatabaseConnection, shared.Node.NodeConnections.PodsDatabaseConnection, shared.Node.NodeConnections.DataAvailabilityDatabaseConnection, latestBatch)

	wgnm.Wait()
}

// Database connectiobnns
//type DatabaseConnections struct {
//	BlockDatabaseConnection            *leveldb.DB
//	TxnDatabaseConnection              *leveldb.DB
//	PodsDatabaseConnection             *leveldb.DB
//	DataAvailabilityDatabaseConnection *leveldb.DB
//	StaticDatabaseConnection           *leveldb.DB
//}
//
//func InitializeDatabaseConnections() DatabaseConnections {
//	var connections DatabaseConnections
//	connections.BlockDatabaseConnection = blocksync.GetBlockDbInstance()
//	connections.TxnDatabaseConnection = blocksync.GetTxDbInstance()
//	connections.PodsDatabaseConnection = blocksync.GetBatchesDbInstance()
//	connections.DataAvailabilityDatabaseConnection = blocksync.GetDaDbInstance()
//	connections.StaticDatabaseConnection = blocksync.GetStaticDbInstance()
//	return connections
//}
