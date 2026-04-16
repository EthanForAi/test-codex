package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"placevr/placevr"
)

func main() {
	node := placevr.Node{
		NodeID:    "node-demo-1",
		Template:  placevr.DefaultTemplateV1(),
		CreatedAt: time.Now(),
	}

	for i := 1; i <= placevr.MaxShotsPerNode; i++ {
		if err := node.AddShot(placevr.Shot{
			ShotID:     fmt.Sprintf("shot-%02d", i),
			FilePath:   fmt.Sprintf("sessions/demo/node-1/%02d.jpg", i),
			CapturedAt: time.Now(),
		}); err != nil {
			log.Fatalf("failed to add shot: %v", err)
		}
	}

	if err := node.Validate(); err != nil {
		log.Fatalf("invalid node: %v", err)
	}

	payload, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		log.Fatalf("marshal failed: %v", err)
	}

	fmt.Println(string(payload))
	fmt.Printf("ready_for_upload=%t\n", node.IsReadyForUpload())
}
