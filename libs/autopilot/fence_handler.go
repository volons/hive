package autopilot

import (
	"fmt"
	"log"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/callback"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

type fenceHandler struct {
	//*events.Listener
	ap *Autopilot

	fence *models.Fence
	state models.FenceState

	goToInProgress *libs.AtomicBool

	done libs.Done
}

func newFenceHandler(ap *Autopilot, fence *models.Fence) *fenceHandler {
	return &fenceHandler{
		//Listener: events.NewListener(),
		ap: ap,

		fence: fence,
		state: models.FenceState{},

		goToInProgress: &libs.AtomicBool{},

		done: libs.NewDone(),
	}
}

//func (h *fenceHandler) start() {
//	positionCh := h.ListenTo(h.ap, "position").(chan models.Position)
//	defer h.StopListening()
//
//	for {
//		select {
//		case pos := <-positionCh:
//			h.checkFence(pos)
//		case <-h.done.WaitCh():
//			return
//		}
//	}
//}

func (h *fenceHandler) checkFence(pos models.Position) {
	if h.goToInProgress.Get() {
		return
	}

	autoRc, targetPos, slowed, outside := h.fence.GetAutoRc(pos, h.ap.manualRc)
	h.ap.SetAutoRc(autoRc)
	h.setFenceStatus(slowed, outside)

	if targetPos != nil {
		h.goTo(*targetPos)
	}
}

func (h *fenceHandler) setFenceStatus(slowed bool, outside bool) {
	var changed bool

	if h.state.Slowed != slowed {
		changed = true
		h.state.Slowed = slowed
	}

	if h.state.Outside != outside {
		changed = true
		h.state.Outside = outside
	}

	if changed {
		h.ap.user.Send(messages.New("fence_state", h.state))
	}
}

func (h *fenceHandler) goTo(pos models.Position) {
	h.goToInProgress.Set(true)
	go func() {
		err := h.ap.GoTo(pos)
		log.Println("GoTo error:", err)
		h.goToInProgress.Set(false)
	}()
}

// GoTo moves the vehicle to the specified position
func (ap *Autopilot) GoTo(pos models.Position) error {
	if !ap.vehicle.Connected() {
		return fmt.Errorf("Vehicle not connected")
	}

	cb := callback.New()
	ap.vehicle.Send(messages.NewRequest("goto", pos, cb))
	_, err := cb.Timeout(time.Minute).Wait()

	return err
}

func (h *fenceHandler) Done() {
	h.done.Done()
}
