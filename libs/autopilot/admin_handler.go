package autopilot

import (
	"errors"

	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

func (ap *Autopilot) handleAdminMessage(msg messages.Message) {
	switch msg.Type {
	case "fence:enable":
		ap.onEnableFence(msg)
	case "fence:disable":
		ap.onDisableFence(msg)
	case "stop":
		ap.onStop()
	}
}

func (ap *Autopilot) onEnableFence(msg messages.Message) {
	err := ap.enableFence()
	msg.Reply(nil, err)
}

func (ap *Autopilot) onDisableFence(msg messages.Message) {
	ap.setFence(nil)
	msg.Reply(nil, nil)
}

func (ap *Autopilot) onStop() {
	ap.done.Done()
}

func (ap *Autopilot) enableFence() error {
	if ap.getFence() != nil {
		return nil
	}

	fence := models.GetFence()
	if fence == nil {
		return errors.New("No fence set")
	}

	ap.setFence(fence)

	return nil
}

func (ap *Autopilot) setFence(fence *models.Fence) {
	if ap.fence != nil {
		ap.fence.Done()
		ap.fence = nil
		ap.autoRc = nil
	}

	if fence != nil {
		ap.fence = newFenceHandler(ap, fence)
		//go ap.fence.start()
	}
}

func (ap *Autopilot) getFence() *fenceHandler {
	return ap.fence
}
