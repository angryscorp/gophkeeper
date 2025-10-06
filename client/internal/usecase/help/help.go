package help

import "gophkeeper/pkg/buildinfo"

type Help struct {
	version   string
	buildTime string
}

func New(version, buildTime string) *Help {
	return &Help{
		version:   version,
		buildTime: buildTime,
	}
}

func (h *Help) Help() string {
	return str + buildinfo.New(h.version, h.buildTime).String() + footer
}

var str = `🔐 GophKeeper — Secure Vault

GophKeeper is a cross-platform CLI client for storing and syncing sensitive data:
  ▸ 💳 Bank card details	
  ▸ 🔑 Logins & passwords
  ▸ 📝 Notes & text
  ▸ 📁 Binary files

🔑 Main actions:
  ▸ Register   — create a new account
  ▸ Login      — sign in to an existing account
  ▸ Sync       — sync data with the server
  ▸ Vault      — view saved items
  ▸ Add Item   — add new data
  ▸ Help       — show this help
  ▸ Quit       — exit the app

🛡 Security:
  ▸ All data is encrypted client-side and transmitted over TLS 1.3
  ▸ The server never stores raw passwords

`

var footer = `
(use ←/esc to return)
`
