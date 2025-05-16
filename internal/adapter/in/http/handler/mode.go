package handler

import (
	"net/http"

	"marketfuck/internal/application/port/in"
)

type ModeHandler struct {
	modeService in.ModeService
}

func NewModeHandler(modeService in.ModeService) *ModeHandler {
	return &ModeHandler{
		modeService: modeService,
	}
}


func (h *ModeHandler) HandleSwitchToTestMode(w http.ResponseWriter, r *http.Request) {
}

func (h *ModeHandler) HandleSwitchToLiveMode(w http.ResponseWriter, r *http.Request) {
}
