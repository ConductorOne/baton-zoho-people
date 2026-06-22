package main

import (
	"github.com/conductorone/baton-sdk/pkg/config"
	cfg "github.com/conductorone/baton-zoho-people/pkg/config"
)

func main() {
	config.Generate("zoho-people", cfg.Config)
}
