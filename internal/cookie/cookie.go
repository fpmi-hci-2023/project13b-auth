package cookie

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/fpmi-hci-2023/project13b-auth/config"
)

// SetCookie sets cookie name and value for TTL seconds
func SetCookie(c *fiber.Ctx, name, value string, ttl int64) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Domain:   config.GlobalConfig.Host,
		MaxAge:   int(ttl),
		Expires:  time.Now().Add(time.Second * time.Duration(ttl)).UTC(),
		Secure:   config.GlobalConfig.SecureCookie,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}

// DeleteCookie deletes cookie by name
func DeleteCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Domain:   config.GlobalConfig.Host,
		MaxAge:   0,
		Expires:  time.Unix(0, 0).UTC(),
		Secure:   config.GlobalConfig.SecureCookie,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}
