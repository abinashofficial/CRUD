package store

import "crud/store/auth"

type Store struct {
	AuthRepo auth.Repository
}
