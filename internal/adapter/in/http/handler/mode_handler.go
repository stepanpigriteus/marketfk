package handler

import (
	"marketfuck/internal/application/port/in"
	"net/http"
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
	err := h.modeService.SwitchToTestMode(r.Context())
	if err != nil {
		http.Error(w, "Failed to switch to test mode: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Switched to TEST mode"))
}

func (h *ModeHandler) HandleSwitchToLiveMode(w http.ResponseWriter, r *http.Request) {
	err := h.modeService.SwitchToLiveMode(r.Context())
	if err != nil {
		http.Error(w, "Failed to switch to live mode: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Switched to LIVE mode"))
}
