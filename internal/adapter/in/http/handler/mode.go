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

// обрабатывает запрос на переключение в тестовый режим
func (h *ModeHandler) HandleSwitchToTestMode(w http.ResponseWriter, r *http.Request) {
}

// обрабатывает запрос на переключение в рабочий режим
func (h *ModeHandler) HandleSwitchToLiveMode(w http.ResponseWriter, r *http.Request) {
}
