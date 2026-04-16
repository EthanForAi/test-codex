package placevr

import (
	"testing"
	"time"
)

func TestDefaultTemplateHasEightSlots(t *testing.T) {
	tpl := DefaultTemplateV1()
	if err := tpl.Validate(); err != nil {
		t.Fatalf("template should be valid: %v", err)
	}
}

func TestNodeRejectsMoreThanEightShots(t *testing.T) {
	node := Node{
		NodeID:    "node-1",
		Template:  DefaultTemplateV1(),
		CreatedAt: time.Now(),
	}

	for i := 0; i < MaxShotsPerNode; i++ {
		err := node.AddShot(Shot{ShotID: "s", FilePath: "/tmp/s.jpg", CapturedAt: time.Now()})
		if err != nil {
			t.Fatalf("shot %d should be accepted: %v", i, err)
		}
	}

	if !node.IsReadyForUpload() {
		t.Fatal("node should be ready after 8 shots")
	}

	if err := node.AddShot(Shot{ShotID: "overflow", FilePath: "/tmp/o.jpg", CapturedAt: time.Now()}); err == nil {
		t.Fatal("expected overflow error")
	}
}
