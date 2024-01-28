package pghera_test

import (
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/TudorHulban/HERA/pgHERA"
	"github.com/gofiber/fiber"
	utils "github.com/gofiber/utils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// TRoute Type to use for insert of model testing.
type TRoute struct {
	Route string
}

// test hera using satisifed interfaces.
func TestWFiber(t *testing.T) {
	h, errCo := pghera.New(pghera.LocalDBInfo, 3, true)
	if assert.Nil(t, errCo) {
		defer h.CloseDBConnection()

		// create table. ignore error, we need the table.
		h.CreateTable(interface{}(&TRoute{}), false)

		var iDDL pghera.IHSQL
		iDDL = h

		app := fiber.New()
		app.Get("/:param", func(c *fiber.Ctx) {
			p := c.Params("param")
			h.L.Debug("p: ", p)
			c.Send(iDDL.InsertModel(&TRoute{
				Route: p,
			}))
		})

		h.L.Debug("starting fiber")
		go app.Listen(3000)

		resp, err := app.Test(httptest.NewRequest("GET", "/"+strconv.FormatInt(time.Now().Unix(), 10), nil))
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	}
}
