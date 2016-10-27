package clusterd

import (
	"path"
	"testing"

	"github.com/rook/rook/pkg/clusterd/inventory"
	"github.com/rook/rook/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestLoadDiscoveredNodes(t *testing.T) {
	etcdClient := &util.MockEtcdClient{}
	mockHandler := newTestServiceLeader()
	raised := make(chan bool)
	mockHandler.unhealthyNode = func(nodes map[string]*UnhealthyNode) {
		assert.Equal(t, 1, len(nodes))
		assert.Equal(t, "23", nodes["23"].ID)
		raised <- true
	}

	context := &Context{EtcdClient: etcdClient}
	context.Services = []*ClusterService{
		&ClusterService{Name: "test", Leader: mockHandler},
	}
	leader := newServicesLeader(context)
	leader.refresher.Start()
	defer leader.refresher.Stop()
	leader.parent = &ClusterMember{isLeader: true}

	etcdClient.SetValue(path.Join(inventory.NodesConfigKey, "23", "publicIp"), "1.2.3.4")
	etcdClient.SetValue(path.Join(inventory.NodesConfigKey, "23", "privateIp"), "10.2.3.4")

	// one unhealthy nodes to discover
	err := leader.discoverUnhealthyNodes()
	<-raised
	assert.Nil(t, err)
}
