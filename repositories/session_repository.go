package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type SessionRepository struct {
	store *session.Store
}

func (s SessionRepository) Set(ctx context.Context, key string, duration string, data interface{}) error {
	c := (ctx).Value("fiberContext").(*fiber.Ctx)

	// Get the session
	sess, err := s.store.Get(c)
	if err != nil {
		return err
	}

	fmt.Println(sess.Reset())

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set data in the session
	sess.Set(key, dataBytes)

	// Save the session
	if err = sess.Save(); err != nil {
		return err
	}

	fmt.Println(sess.Keys())

	return nil
}

func (s SessionRepository) Get(ctx context.Context, key string) ([]byte, error) {
	// Get the Fiber context from the GraphQL context
	c := (ctx).Value("fiberContext").(*fiber.Ctx)

	// Get the session
	sess, err := s.store.Get(c)
	if err != nil {
		return nil, err
	}

	// Get value
	name := sess.Get(key)

	// Get all Keys
	keys := sess.Keys()

	fmt.Println(keys)

	bodyBytes, _ := json.Marshal(name)

	return bodyBytes, nil
}

func (s SessionRepository) Del(ctx context.Context, key string) error {
	// Get the Fiber context from the GraphQL context
	c := (ctx).Value("fiberContext").(*fiber.Ctx)

	// Get the session
	sess, err := s.store.Get(c)
	if err != nil {
		return err
	}

	// Retrieve data from the session
	sess.Delete(key)

	// Save the session
	if err = sess.Save(); err != nil {
		return err
	}

	return nil
}

type SessionRepositoryInterface interface {
	Set(ctx context.Context, key string, duration string, data interface{}) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

func NewSessionRepository(store *session.Store) SessionRepositoryInterface {
	return &SessionRepository{store: store}
}
