package role

type (
	RegisterInputRole struct {
		Name string `json:"name" validate:"required"`
	}

	UpdateInputRole struct {
		Name    string `json:"name" validate:"required"`
		Version int64  `json:"version" validate:"required"`
	}

	SoftDeleteInputRole struct {
		Version int64 `json:"version" validate:"required"`
	}
)
