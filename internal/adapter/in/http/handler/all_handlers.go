package handler

import (
	"marketfuck/internal/application/port"
	"marketfuck/internal/application/port/in"
)

type AllHandlers struct {
	Health *HealthHandler
	Mode   *ModeHandler
	Price  *PriceHandler
}

func NewAllHandlers(healthService in.HealthService, modeService in.ModeService, priceService in.PriceService, logger port.Logger) *AllHandlers {
	return &AllHandlers{
		Health: NewHealthHandler(healthService),
		Mode:   NewModeHandler(modeService),
		Price:  NewPriceHandler(priceService, logger),
	}
}
