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

var ctx = context.Background()
var client, _ = container.NewClusterManagerClient(ctx, option.WithCredentialsFile(os.Getenv("GCKEY_FILE_PATH")))

func getPools() []string {
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
	var names = make([]string, len(nodePools))
	for index, nodePool := range nodePools {
		names[index] = nodePool.Name
	}
	return names
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
	nodePoolNames := getPools()
	// You can only resize one nodePool size at a time
	// and have to wait for GKE done resizing
	setNodePoolSize(nodePoolNames[0], 0)
	// TODO: waiting for resizing job done and set another nodePool size
}
