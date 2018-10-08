package eth

import (
	"reflect"
	"testing"

	"bitbucket.org/cpchain/chain/common"
	"bitbucket.org/cpchain/chain/p2p"
)

// launch the chain
// new a committee_network_handler
// build the network.

func TestNewBasicCommitteeNetworkHandler(t *testing.T) {
	type args struct {
		peers           *peerSet
		epochLength     uint64
		ownAddress      common.Address
		contractAddress common.Address
		server          *p2p.Server
	}
	tests := []struct {
		name    string
		args    args
		want    *BasicCommitteeNetworkHandler
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBasicCommitteeNetworkHandler(tt.args.peers, tt.args.epochLength, tt.args.ownAddress, tt.args.contractAddress, tt.args.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBasicCommitteeNetworkHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBasicCommitteeNetworkHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}