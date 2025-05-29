package models

import "time"

type UserAuth struct {
	id           string
	email        string
	passwordHash string
	name         string
	roleId       string
	createdAt    time.Time
}
