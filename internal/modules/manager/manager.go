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
	ReplacePair(kv *models.KV) error
	DeletePairByKey(key string) error
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

// UpdateValue replaces KV by key in storage.
func (m *Manager) UpdateValue(kv *models.KV) error {
	err := m.storage.ReplacePair(kv)
	if err != nil {
		if errors.Is(err, errs.ErrNoKeys) {
			slog.Error("Error while updating kv. Key not found", "key:", kv.Key)
		}
		return errs.Wrap(err)
	}

	slog.Info("Pair successfully updated. Key:", "key:", kv.Key)
	return nil
}

// DeleteValue deletes KV by key in storage.
func (m *Manager) DeleteValue(key string) error {
	err := m.storage.DeletePairByKey(key)
	if err != nil {
		if errors.Is(err, errs.ErrNoKeys) {
			slog.Error("Error while updating kv. Key not found", "key:", key)
		}
		return errs.Wrap(err)
	}

	slog.Info("Pair successfully deleted. Key:", "key:", key)
	return nil
}
