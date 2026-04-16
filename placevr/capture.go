package placevr

import (
	"errors"
	"fmt"
	"time"
)

const (
	// MaxShotsPerNode enforces the UX constraint for MVP capture flow.
	MaxShotsPerNode = 8
)

var (
	ErrTooManyShots       = fmt.Errorf("node has more than %d shots", MaxShotsPerNode)
	ErrTemplateSize       = fmt.Errorf("capture template must contain exactly %d guide slots", MaxShotsPerNode)
	ErrTemplateOrder      = errors.New("capture template order must be continuous starting from 1")
	ErrDuplicateGuideSlot = errors.New("capture template contains duplicate order slot")
)

// GuideSlot defines one guided capture position in a fixed template.
type GuideSlot struct {
	Order int     `json:"order"`
	Yaw   float64 `json:"yaw"`
	Pitch float64 `json:"pitch"`
}

// CaptureTemplate is fixed in MVP to improve completion rate and upload predictability.
type CaptureTemplate struct {
	ID    string      `json:"id"`
	Slots []GuideSlot `json:"slots"`
}

func (t CaptureTemplate) Validate() error {
	if len(t.Slots) != MaxShotsPerNode {
		return ErrTemplateSize
	}

	seen := make(map[int]struct{}, MaxShotsPerNode)
	for _, slot := range t.Slots {
		if _, ok := seen[slot.Order]; ok {
			return ErrDuplicateGuideSlot
		}
		seen[slot.Order] = struct{}{}
	}

	for i := 1; i <= MaxShotsPerNode; i++ {
		if _, ok := seen[i]; !ok {
			return ErrTemplateOrder
		}
	}

	return nil
}

// DefaultTemplateV1 gives 8 directions around the user with slight pitch variation.
func DefaultTemplateV1() CaptureTemplate {
	return CaptureTemplate{
		ID: "indoor_8shot_v1",
		Slots: []GuideSlot{
			{Order: 1, Yaw: 0, Pitch: 0},
			{Order: 2, Yaw: 45, Pitch: 0},
			{Order: 3, Yaw: 90, Pitch: 0},
			{Order: 4, Yaw: 135, Pitch: 0},
			{Order: 5, Yaw: 180, Pitch: 0},
			{Order: 6, Yaw: 225, Pitch: 0},
			{Order: 7, Yaw: 270, Pitch: 0},
			{Order: 8, Yaw: 315, Pitch: -8},
		},
	}
}

// Shot is one original image captured on the phone.
type Shot struct {
	ShotID     string    `json:"shotId"`
	FilePath   string    `json:"filePath"`
	CapturedAt time.Time `json:"capturedAt"`
}

// Node is the minimum cloud stitching unit.
type Node struct {
	NodeID    string          `json:"nodeId"`
	Template  CaptureTemplate `json:"template"`
	Shots     []Shot          `json:"shots"`
	CreatedAt time.Time       `json:"createdAt"`
}

func (n Node) Validate() error {
	if err := n.Template.Validate(); err != nil {
		return err
	}
	if len(n.Shots) > MaxShotsPerNode {
		return ErrTooManyShots
	}
	return nil
}

func (n Node) IsReadyForUpload() bool {
	return len(n.Shots) == MaxShotsPerNode
}

func (n *Node) AddShot(shot Shot) error {
	if len(n.Shots) >= MaxShotsPerNode {
		return ErrTooManyShots
	}
	n.Shots = append(n.Shots, shot)
	return nil
}
