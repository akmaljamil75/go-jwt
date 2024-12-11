package user

type (
	RegisterInputUser struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		RoleID   uint   `json:"role_id" validate:"required"`
	}

	SoftDeleteInputUser struct {
		Version int64 `json:"version" validate:"required"`
		ID      uint  `json:"id" validate:"required"`
	}

	UpdateInputUser struct {
		Username string `json:"username"`
		Password string `json:"password" `
		RoleID   uint   `json:"role_id"`
		Version  int64  `json:"version" validate:"required"`
	}
)
