package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

type formatArg struct {
	Format string
}

const (
	FormatArgJson   = "json"
	FormatArgPretty = "pretty"
)

func (f *formatArg) UnmarshalText(text []byte) error {
	f.Format = string(text)

	if f.Format != FormatArgJson && f.Format != FormatArgPretty {
		return fmt.Errorf("invalid format: %s", f.Format)
	}
	return nil
}

type Args struct {
	Format  formatArg `arg:"--format" help:"Output format (json, pretty)" default:"json"`
	Address string    `arg:"-a,env:ADDRESS" help:"Address to listen on" default:"*"`
	Port    int       `arg:"-p,env:PORT" help:"Port to listen on" default:"8080"`
	// Postgres
	PostgresHost     string `arg:"--postgres-host,required,env:POSTGRES_HOST" help:"Postgres host"`
	PostgresPort     uint16 `arg:"--postgres-port,env:POSTGRES_PORT" help:"Postgres port" default:"5432"`
	PostgresDatabase string `arg:"--postgres-database,required,env:POSTGRES_DATABASE" help:"Postgres database"`
	PostgresUser     string `arg:"--postgres-user,required,env:POSTGRES_USER" help:"Postgres user"`
	PostgresPassword string `arg:"--postgres-password,required,env:POSTGRES_PASSWORD" help:"Postgres password"`
	PostgresTls      bool   `arg:"--postgres-tls,env:POSTGRES_TLS" help:"Use TLS for Postgres" default:"false"`
}

var parsedArgs bool
var globalArgs Args

func parseArgs() Args {
	parsedArgs = true
	arg.MustParse(&globalArgs)
	return globalArgs
}

func GetArgs() Args {
	if !parsedArgs {
		parseArgs()
	}
	return globalArgs
}
