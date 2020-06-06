package pghera_test

import (
	"testing"

	"github.com/TudorHulban/HERA/pgHERA"
	"github.com/gofiber/fiber"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// TRoute Type to use for insert of model testing.
type TRoute struct {
	Route string
}

func TestWFiber(t *testing.T) {
	var errCo error
	h, errCo = pghera.New(pghera.LocalDBInfo, 3, true)

	if assert.Nil(t, errCo) {
		app := fiber.New()

		app.Get("/", func(c *fiber.Ctx) {
			c.Send("Hello, World!")
		})

		app.Listen(3000)
	}
}
