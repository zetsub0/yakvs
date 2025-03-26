package manager

import (
	"errors"
	"log/slog"

	"github.com/zetsub0/yakvs/internal/models"
	"github.com/zetsub0/yakvs/pkg/errs"
)

// Storage implement kv database
type Storage interface {
	GetPairByKey(key string) (*models.KV, error)
	SetPair(kv *models.KV) error
}

// Manager ...
type Manager struct {
	storage Storage
}

// New ...
func New(storage Storage) *Manager {
	return &Manager{
		storage: storage,
	}
}

// CreateValue send create KV request to storage.
func (m *Manager) CreateValue(kv *models.KV) error {
	err := m.storage.SetPair(kv)
	if err != nil {
		if errors.Is(err, errs.ErrKeyExists) {
			slog.Error("Error while creating kv. Key already exists", "key:", kv.Key)
		}
		return errs.Wrap(err)
	}

	slog.Info("Pair successfully inserted. Key:", "key:", kv.Key)
	return nil
}

// GetValue receives KV by key if it exist in storage.
func (m *Manager) GetValue(key string) (*models.KV, error) {
	pair, err := m.storage.GetPairByKey(key)
	if err != nil {
		slog.Error(err.Error())
		return nil, errs.Wrap(err)
	}

	return pair, nil
}
