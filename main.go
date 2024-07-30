package main

import (
	"context"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var dnsServers = map[string]string{
	"cloudflare": "1.1.1.1:53",
	"opendns":    "208.67.222.222:53",
	"google":     "8.8.8.8:53",
	"yandex":     "77.88.8.8:53",
}

func main() {
	app := fiber.New(fiber.Config{})
	app.Use(cors.New())
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   5000 * time.Millisecond,
		CacheControl: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/a/:dns/:domain", getARecord)
	app.Get("/aaaa/:dns/:domain", getAAAARecord)
	app.Get("/cname/:dns/:domain", getCNAMERecord)
	app.Get("/mx/:dns/:domain", getMXRecord)
	app.Get("/ns/:dns/:domain", getNSRecord)
	app.Get("/ptr/:dns/:ip", getPTRRecord)
	app.Get("/txt/:dns/:domain", getTXTRecord)
	//SOA

	app.Get("/all/:dns/:domain", getALLforDomain)

	app.Listen(":3000")
}

func getARecord(c *fiber.Ctx) error {
	dns := c.Params("dns")
	domain := c.Params("domain")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	ipv4s, err := resolver.LookupIP(context.Background(), "ip4", domain)
	if err != nil {
		println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "No records found")
	}

	return c.JSON(
		fiber.Map{
			"ipv4": ipv4s,
		},
	)
}

func getAAAARecord(c *fiber.Ctx) error {
	//grab parabeter
	dns := c.Params("dns")
	domain := c.Params("domain")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	ipv6s, err := resolver.LookupIP(context.Background(), "ip6", domain)
	if err != nil {
		println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "No records found")
	}
	return c.JSON(
		fiber.Map{
			"ipv6": ipv6s,
		},
	)
}

func getCNAMERecord(c *fiber.Ctx) error {
	dns := c.Params("dns")
	domain := c.Params("domain")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	cname, err := resolver.LookupCNAME(context.Background(), domain)
	if err != nil {
		println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "No records found")
	}
	return c.JSON(
		fiber.Map{
			"cname": cname,
		},
	)
}

func getMXRecord(c *fiber.Ctx) error {
	dns := c.Params("dns")
	domain := c.Params("domain")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	mx, err := resolver.LookupMX(context.Background(), domain)
	if err != nil {
		println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "No records found")
	}

	return c.JSON(
		fiber.Map{
			"mx": mx,
		},
	)
}

func getNSRecord(c *fiber.Ctx) error {
	dns := c.Params("dns")
	domain := c.Params("domain")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	ns, err := resolver.LookupNS(context.Background(), domain)
	if err != nil {
		println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "No records found")
	}
	return c.JSON(
		fiber.Map{
			"ns": ns,
		},
	)
}

func getPTRRecord(c *fiber.Ctx) error {
	dns := c.Params("dns")
	ip := c.Params("ip")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	ptr, err := resolver.LookupAddr(context.Background(), ip)
	if err != nil {
		println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "No records found")
	}
	return c.JSON(
		fiber.Map{
			"ptr": ptr,
		},
	)
}

func getTXTRecord(c *fiber.Ctx) error {
	dns := c.Params("dns")
	domain := c.Params("domain")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	txt, err := resolver.LookupTXT(context.Background(), domain)
	if err != nil {
		println(err.Error())
		return fiber.NewError(fiber.StatusNotFound, "No records found")
	}
	return c.JSON(
		fiber.Map{
			"txt": txt,
		},
	)
}

func getALLforDomain(c *fiber.Ctx) error {
	dns := c.Params("dns")
	domain := c.Params("domain")
	dnsIP, ok := dnsServers[dns]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "DNS server not found")
	}
	resolver := makeResolver(dnsIP)
	ipv4s, _ := resolver.LookupIP(context.Background(), "ip4", domain)
	ipv6s, _ := resolver.LookupIP(context.Background(), "ip6", domain)
	cname, _ := resolver.LookupCNAME(context.Background(), domain)
	mx, _ := resolver.LookupMX(context.Background(), domain)
	ns, _ := resolver.LookupNS(context.Background(), domain)
	txt, _ := resolver.LookupTXT(context.Background(), domain)

	return c.JSON(
		fiber.Map{
			"ipv4":  ipv4s,
			"ipv6":  ipv6s,
			"cname": cname,
			"mx":    mx,
			"ns":    ns,
			"txt":   txt,
		},
	)
}
func makeResolver(dnsIP string) net.Resolver {
	return net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, "tcp", dnsIP)
		},
	}
}
