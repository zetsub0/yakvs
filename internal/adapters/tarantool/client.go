package tarantool

import (
	"context"
	"errors"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/zetsub0/yakvs/internal/config"
	"github.com/zetsub0/yakvs/internal/models"
	"github.com/zetsub0/yakvs/pkg/errs"
)

// space is constant for space name.
const space = "kv"

// Tarantool is implementation for tarantool client
type Tarantool struct {
	client *tarantool.Connection
}

// New returns Tarantool client
func New(ctx context.Context, cfg config.Tarantool) *Tarantool {
	dialer := &tarantool.NetDialer{
		Address:  cfg.Host,
		User:     cfg.User,
		Password: cfg.Password,
	}

	opts := tarantool.Opts{
		Timeout: cfg.Timeout,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		panic(err)
	}

	_, err = conn.Do(tarantool.NewPingRequest()).Get()
	if err != nil {
		panic(err)
	}

	return &Tarantool{client: conn}
}

// SetPair inserts nev kv to vault.
func (t *Tarantool) SetPair(kv *models.KV) error {

	res, err := t.getPairByKey(kv.Key)
	if err != nil && !errors.Is(err, errs.ErrNoKeys) {
		return errs.Wrapf(err, "error got while validating key")
	}

	if res != nil {
		return errs.Wrap(errs.ErrKeyExists)
	}

	tuple := []any{kv.Key, kv.Value}

	req := tarantool.NewInsertRequest(space).Tuple(tuple)

	_, err = t.client.Do(req).Get()
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

// GetPairByKey returns kv from vault by inputted key.
func (t *Tarantool) GetPairByKey(key string) (*models.KV, error) {
	return t.getPairByKey(key)
}

// getPairByKey is internal method that returns kv from vault by inputted key.
func (t *Tarantool) getPairByKey(key string) (*models.KV, error) {

	res, err := t.client.Do(tarantool.NewSelectRequest(space).Key([]any{key})).Get()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if len(res) == 0 {
		return nil, errs.Wrap(errs.ErrNoKeys)
	}

	innerSlice, ok := res[0].([]any)
	if !ok {
		return nil, errs.Wrap(errors.New("failed to convert data from tarantool"))
	}

	if len(innerSlice) < 2 {
		return nil, errs.Wrap(errors.New("kv malformed"))
	}

	k, ok := innerSlice[0].(string)
	if !ok {
		return nil, errs.Wrap(errors.New("key malformed"))
	}
	v, ok := innerSlice[1].([]byte)
	if !ok {
		return nil, errs.Wrap(errors.New("value malformed"))
	}

	return &models.KV{
		Key:   k,
		Value: v,
	}, nil
}
