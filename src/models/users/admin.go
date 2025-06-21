package usermodel

type AdminUser struct {
	Username   string `json:"username" validate:"required"`
	Role       string `json:"role" validate:"required"`
	HashedPW   string `json:"password_hash" validate:"required"`
	CreatedAt  int64  `json:"created_at" validate:"required"`
	ModifiedAt int64  `json:"modified_at" validate:"required"`
}
