package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"golang.org/x/crypto/sha3"
)

var (
	blockchainID = ""
	networkID    = ""
)

type APIServer struct {
}

func (s APIServer) AccountBalance(ctx context.Context, r *types.AccountBalanceRequest) (*types.AccountBalanceResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) AccountCoins(ctx context.Context, r *types.AccountCoinsRequest) (*types.AccountCoinsResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) Block(ctx context.Context, r *types.BlockRequest) (*types.BlockResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) BlockTransaction(ctx context.Context, r *types.BlockTransactionRequest) (*types.BlockTransactionResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) Call(ctx context.Context, r *types.CallRequest) (*types.CallResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) ConstructionCombine(ctx context.Context, r *types.ConstructionCombineRequest) (*types.ConstructionCombineResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) ConstructionDerive(ctx context.Context, r *types.ConstructionDeriveRequest) (*types.ConstructionDeriveResponse, *types.Error) {
	key := append(r.PublicKey.Bytes, 0)
	hash := sha3.Sum256(key)
	return &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: hex.EncodeToString(hash[:]),
		},
	}, nil
}

func (s APIServer) ConstructionHash(ctx context.Context, r *types.ConstructionHashRequest) (*types.TransactionIdentifierResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) ConstructionMetadata(ctx context.Context, r *types.ConstructionMetadataRequest) (*types.ConstructionMetadataResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) ConstructionParse(ctx context.Context, r *types.ConstructionParseRequest) (*types.ConstructionParseResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) ConstructionPayloads(ctx context.Context, r *types.ConstructionPayloadsRequest) (*types.ConstructionPayloadsResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) ConstructionPreprocess(ctx context.Context, r *types.ConstructionPreprocessRequest) (*types.ConstructionPreprocessResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) ConstructionSubmit(ctx context.Context, r *types.ConstructionSubmitRequest) (*types.TransactionIdentifierResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) Mempool(ctx context.Context, r *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) MempoolTransaction(ctx context.Context, r *types.MempoolTransactionRequest) (*types.MempoolTransactionResponse, *types.Error) {
	return nil, nil
}

func (s APIServer) NetworkList(ctx context.Context, r *types.MetadataRequest) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{{
			Blockchain: blockchainID,
			Network:    networkID,
		}},
	}, nil
}

func (s APIServer) NetworkOptions(ctx context.Context, r *types.NetworkRequest) (*types.NetworkOptionsResponse, *types.Error) {
	return &types.NetworkOptionsResponse{
		Allow: &types.Allow{
			Errors: []*types.Error{
				{Code: 101, Message: "some error", Retriable: false},
			},
			HistoricalBalanceLookup: true,
			OperationStatuses: []*types.OperationStatus{
				{Status: "FAILURE", Successful: false},
				{Status: "SUCCESS", Successful: true},
			},
			OperationTypes: []string{"foo"},
		},
		Version: &types.Version{
			MiddlewareVersion: types.String("0.0.1"),
			NodeVersion:       "0.0.1",
			RosettaVersion:    "0.0.1",
		},
	}, nil
}

func (s APIServer) NetworkStatus(ctx context.Context, r *types.NetworkRequest) (*types.NetworkStatusResponse, *types.Error) {
	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Hash:  "latest",
			Index: 1,
		},
		CurrentBlockTimestamp: time.Now().UnixMilli(),
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Hash:  "genesis",
			Index: 0,
		},
	}, nil
}

func main() {
	blockchainID = os.Getenv("DUMMY_BLOCKCHAIN")
	if blockchainID == "" {
		log.Fatalf("Missing DUMMY_BLOCKCHAIN environment value")
	}
	networkID = os.Getenv("DUMMY_NETWORK")
	if networkID == "" {
		log.Fatalf("Missing DUMMY_NETWORK environment value")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	networks := []*types.NetworkIdentifier{{
		Blockchain: blockchainID,
		Network:    networkID,
	}}
	asserter, err := asserter.NewServer(
		[]string{"foo"},
		true,
		networks,
		[]string{},
		false,
		"",
	)
	if err != nil {
		log.Fatalf("Failed to instantiate the Rosetta asserter: %s", err)
	}
	api := APIServer{}
	router := server.NewRouter(
		server.NewAccountAPIController(api, asserter),
		server.NewBlockAPIController(api, asserter),
		server.NewCallAPIController(api, asserter),
		server.NewConstructionAPIController(api, asserter),
		server.NewMempoolAPIController(api, asserter),
		server.NewNetworkAPIController(api, asserter),
	)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	log.Printf("Starting Rosetta Server on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Rosetta HTTP Server failed: %s", err)
	}
}
