package in

import "context"

type ModeService interface {
	// в тестовый режим 
	SwitchToTestMode(ctx context.Context) error
	//  в рабочий режим 
	SwitchToLiveMode(ctx context.Context) error
}
