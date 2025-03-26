package manager

import (
	"errors"
	"log/slog"

	"github.com/zetsub0/yakvs/internal/models"
	"github.com/zetsub0/yakvs/pkg/errs"
)

type Storage interface {
	GetPairByKey(key string) (*models.KV, error)
	SetPair(kv *models.KV) error
}

type Manager struct {
	storage Storage
}

func New(storage Storage) *Manager {
	return &Manager{
		storage: storage,
	}
}

func (m *Manager) CreateValue(kv *models.KV) error {
	err := m.storage.SetPair(kv)
	if err != nil {
		if errors.Is(err, errs.ErrKeyExists) {
			slog.Error("Error while creating kv. Key already exists", "key:", kv.Key)
			return errs.Wrap(errs.ErrKeyExists)
		} else {
			return errs.Wrap(err)
		}
	}

	slog.Info("Pair successfully inserted. Key:", "key:", kv.Key)
	return nil
}
