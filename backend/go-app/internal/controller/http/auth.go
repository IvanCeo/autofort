//go:build ignore

package http

import (
	"fmt"
	"os"
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

type Credentials struct {
	login    string `json:"login"`
	password string `json:"password"`
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	req := Credentials{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).SendString(fmt.Sprintf(`{"error":"%v"}`, err))
	}

	// как это должно делаться
	//err = bcrypt.CompareHashAndPassword([]byte(passhash), []byte(req.password))
	// так не делается потому что нужна ручка + лендинг на регистрацию + надо таблицу в базе заводить + архитектурно чтобы контролеры в базу не ходили

	// рельность :(
	if req.login == os.Getenv("ADMIN_USER") && req.password == os.Getenv("ADMIN_PASS") {
		token, err := h.server.GenerateToken(req.login, 24*time.Hour)
		if err != nil {
			return 500
		}
	}

}
