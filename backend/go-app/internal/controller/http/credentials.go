package http

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) Refresh(c *fiber.Ctx) error {
	// проверяет из httpOnly кук refresh, при успехе выдает
	// новые оба, старый refresh затирает из кук
	refresh := c.Cookies("refresh")
	if refresh == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// он берется из кук, только туда надо его полоить при signIn
	userIDstr := c.Cookies("userID")
	if userIDstr == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	tokens, err := h.server.Refresh(c.Context(), &userID, refresh)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     refresh,
		Value:    tokens.GetRefresh(),
		HTTPOnly: true,
	})

	return c.Status(200).SendString() // TODO: идиоматично отправить json
}

// type Cookie struct {
// 	Name        string    `json:"name"`
// 	Value       string    `json:"value"`
// 	Path        string    `json:"path"`
// 	Domain      string    `json:"domain"`
// 	MaxAge      int       `json:"max_age"`
// 	Expires     time.Time `json:"expires"`
// 	Secure      bool      `json:"secure"`
// 	HTTPOnly    bool      `json:"http_only"`
// 	SameSite    string    `json:"same_site"`
// 	SessionOnly bool      `json:"session_only"`
// }
