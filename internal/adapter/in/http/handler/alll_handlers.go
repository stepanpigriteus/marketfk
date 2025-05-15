package handler

type AllHandlers struct {
	Health *HealthHandler
	Mode *ModeHandler
	Price *PriceHandler
}
