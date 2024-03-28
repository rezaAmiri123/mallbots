//go:build integration || database

package adapters

import(
_	"github.com/docker/go-connections/nat"
_	 "github.com/jackc/pgx/v4/stdlib"
_	"github.com/pressly/goose/v3"
_	"github.com/stretchr/testify/mock"
_	"github.com/stretchr/testify/suite"
_	"github.com/testcontainers/testcontainers-go"
_	"github.com/testcontainers/testcontainers-go/wait"

)
