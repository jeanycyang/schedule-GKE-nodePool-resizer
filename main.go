package main

import (
	"context"
	"fmt"
	"log"
	"os"

	container "cloud.google.com/go/container/apiv1"
	option "google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

// CommonFields are basic settings on GCP
type CommonFields struct {
	ProjectId string
	Zone      string
	ClusterId string
}

var commonFields = &CommonFields{
	ProjectId: os.Getenv("PROJECT_ID"),
	Zone:      os.Getenv("ZONE"),
	ClusterId: os.Getenv("CLUSTER_ID"),
}

// NodePoolStatus contains useful node pool status
type NodePoolStatus struct {
	Name             string
	InitialNodeCount int32
}

// NodePoolDesiredStatus:
type NodePoolDesiredStatus struct {
	Name             string
	DesiredNodeCount int32
}

var ctx = context.Background()
var client, _ = container.NewClusterManagerClient(ctx, option.WithCredentialsFile(os.Getenv("GCKEY_FILE_PATH")))

func getNodePoolStatuses() []NodePoolStatus {
	req := &containerpb.ListNodePoolsRequest{
		ProjectId: commonFields.ProjectId,
		Zone:      commonFields.Zone,
		ClusterId: commonFields.ClusterId,
	}
	resp, err := client.ListNodePools(ctx, req)
	if err != nil {
		log.Fatal(err.Error())
	}
	nodePools := resp.NodePools
	var statuses = make([]NodePoolStatus, len(nodePools))
	for index, nodePool := range nodePools {
		statuses[index] = NodePoolStatus{
			Name:             nodePool.Name,
			InitialNodeCount: nodePool.InitialNodeCount,
		}
	}
	return statuses
}

func setNodePoolSize(nodePoolId string, nodeCount int32) {
	req := &containerpb.SetNodePoolSizeRequest{
		ProjectId:  commonFields.ProjectId,
		Zone:       commonFields.Zone,
		ClusterId:  commonFields.ClusterId,
		NodePoolId: nodePoolId,
		NodeCount:  nodeCount,
	}

	resp, err := client.SetNodePoolSize(ctx, req)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%+v\n", resp)
}

func main() {
	nodePoolStatuses := getNodePoolStatuses()
	nodePoolDesiredStatuses := [1]NodePoolDesiredStatus{}
	// @TODO: desired statuses should be an argument
	nodePoolDesiredStatuses[0] = NodePoolDesiredStatus{
		Name:             "pool-1",
		DesiredNodeCount: 1,
	}
	// You can only resize one nodePool size at a time
	// and have to wait for GKE done resizing
	for _, desiredStatus := range nodePoolDesiredStatuses {
		for _, nodePoolStatus := range nodePoolStatuses {
			if desiredStatus.Name == nodePoolStatus.Name &&
				desiredStatus.DesiredNodeCount != nodePoolStatus.InitialNodeCount {
				setNodePoolSize(desiredStatus.Name, desiredStatus.DesiredNodeCount)
			}
		}
	}
	// TODO: waiting for resizing job done and set another nodePool size
}
