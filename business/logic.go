package business

import "github.com/chrisvdg/gotiny/backend"

// NewLogic creates a new Logic instance
func NewLogic() (*Logic, error) {

	return nil, nil
}

// Logic contains a stateful set of business logic
type Logic struct {
	backend backend.Backend
}
