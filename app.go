package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/yaegashi/azbill/store"
)

const (
	AppClientID    = "4a034c56-da44-48ce-90db-039a408974bd"
	AppTenantID    = "common"
	ProgressScale  = 100
	ProgressColumn = 50
	AuthStoreEnv   = "AZBILL_AUTH_STORE"
)

type App struct {
	Authorizer autorest.Authorizer
	Writer     io.WriteCloser
	CSVWriter  *csv.Writer
	Output     string
	Auth       string
	AuthStore  string
	Client     string
	Tenant     string
	Format     string
	Flatten    bool
	Pretty     bool
	IsStdout   bool
	Quiet      bool
	Lines      int
	Column     int
}

func (app *App) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "azbill",
		Short:             "Download Azure billing information",
		PersistentPreRunE: app.PersistentPreRunE,
		SilenceUsage:      true,
	}
	cmd.PersistentFlags().StringVarP(&app.Client, "client", "", AppClientID, "Azure app client")
	cmd.PersistentFlags().StringVarP(&app.Tenant, "tenant", "", AppTenantID, "Azure tenant")
	cmd.PersistentFlags().StringVarP(&app.Auth, "auth", "", "dev", "Auth source [dev,env,file=PATH,cli]")
	cmd.PersistentFlags().StringVarP(&app.AuthStore, "auth-store", "", "", "Auth persistent token store")
	cmd.PersistentFlags().StringVarP(&app.Format, "format", "", "csv", "Output format [csv,json,flatten,pretty]")
	cmd.PersistentFlags().StringVarP(&app.Output, "output", "o", "", "Output file")
	cmd.PersistentFlags().BoolVarP(&app.Quiet, "quiet", "q", false, "Quiet")
	return cmd
}

func (app *App) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	app.IsStdout = app.Output == "" || app.Output == "-"
	if app.IsStdout && isatty.IsTerminal(os.Stdout.Fd()) {
		app.Quiet = true
	}

	for _, f := range strings.Split(strings.ToLower(app.Format), ",") {
		switch f {
		case "json":
			app.Format = "json"
		case "csv":
			app.Format = "csv"
		case "flatten":
			app.Format = "json"
			app.Flatten = true
		case "pretty":
			app.Format = "json"
			app.Pretty = true
		default:
			return fmt.Errorf("Unknown format: %s", f)
		}
	}

	var authorizer autorest.Authorizer
	var token *adal.ServicePrincipalToken
	var err error

	if app.AuthStore == "" {
		app.AuthStore = os.Getenv(AuthStoreEnv)
	}
	if app.AuthStore != "" {
		err = func() error {
			b, err := store.Load(app.AuthStore)
			if err != nil {
				return nil
			}
			err = json.Unmarshal(b, &token)
			if err != nil {
				return nil
			}
			save := false
			token.SetRefreshCallbacks([]adal.TokenRefreshCallback{func(adal.Token) error { save = true; return nil }})
			err = token.EnsureFresh()
			if err != nil {
				return nil
			}
			if save {
				b, err = json.Marshal(token)
				if err != nil {
					return err
				}
				err = store.Save(app.AuthStore, b, 0600)
				if err != nil {
					return err
				}
			}
			authorizer = autorest.NewBearerAuthorizer(token)
			return nil
		}()
		if err != nil {
			return err
		}
	}

	if authorizer == nil && app.Auth == "env" {
		authorizer, err = auth.NewAuthorizerFromEnvironment()
		if err != nil {
			fmt.Fprintf(os.Stderr, "auth-env: %s\n", err)
		}
	}

	if authorizer == nil && strings.HasPrefix(app.Auth, "file=") {
		os.Setenv("AZURE_AUTH_LOCATION", app.Auth[5:])
		authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
		if err != nil {
			fmt.Fprintf(os.Stderr, "auth-file: %s\n", err)
		}
	}

	if authorizer == nil && app.Auth == "cli" {
		authorizer, err = auth.NewAuthorizerFromCLI()
		if err != nil {
			fmt.Fprintf(os.Stderr, "auth-cli: %s\n", err)
		}
	}

	if authorizer == nil && app.Auth == "dev" {
		deviceConfig := auth.NewDeviceFlowConfig(app.Client, app.Tenant)
		token, err := deviceConfig.ServicePrincipalToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "auth-dev: %s\n", err)
		}
		if app.AuthStore != "" {
			b, err := json.Marshal(token)
			if err != nil {
				return err
			}
			err = store.Save(app.AuthStore, b, 0600)
			if err != nil {
				return err
			}
		}
		authorizer = autorest.NewBearerAuthorizer(token)
	}

	if authorizer == nil {
		return fmt.Errorf("Failed to configure authorizer")
	}

	app.Authorizer = authorizer

	return nil
}

func (app *App) Open() error {
	format := app.Format
	if format == "json" {
		if app.Flatten {
			format += ",flatten"
		}
		if app.Pretty {
			format += ",pretty"
		}
	}
	if app.IsStdout {
		app.Writer = os.Stdout
		app.Logf("Writing to stdout in %s", format)
	} else {
		w, err := os.Create(app.Output)
		if err != nil {
			return err
		}
		app.Writer = w
		app.Logf("Writing to %q in %s", app.Output, format)
	}
	if app.Format == "csv" {
		app.CSVWriter = csv.NewWriter(app.Writer)
		app.CSVWriter.UseCRLF = true
		app.Writer.Write([]byte{0xef, 0xbb, 0xbf})
	}
	return nil
}

func (app *App) Close() {
	if app.Format == "csv" {
		app.CSVWriter.Flush()
	}
	if !app.IsStdout {
		app.Writer.Close()
	}
	app.Progress(0)
	app.Log("Done")
}

func (app *App) Progress(n int) {
	if n == 0 {
		if app.Column > 0 {
			spaces := ""
			for i := app.Column; i < ProgressColumn; i++ {
				spaces += " "
			}
			app.Msgf("%s %8d lines\n", spaces, app.Lines)
		}
		return
	}
	lines := app.Lines + n
	blocks := app.Lines / ProgressScale
	n = (lines / ProgressScale) - blocks
	app.Lines = lines
	for i := 0; i < n; i++ {
		app.Msg(".")
		app.Column++
		blocks++
		if app.Column >= ProgressColumn {
			app.Msgf(" %8d lines\n", blocks*ProgressScale)
			app.Column = 0
		}
	}
}

func (app *App) JSONMarshal(v interface{}) error {
	enc := json.NewEncoder(app.Writer)
	if app.Pretty {
		enc.SetIndent("", "  ")
	}
	err := enc.Encode(v)
	if err != nil {
		return err
	}
	app.Progress(1)
	return nil
}

func (app *App) CSVMarshal(row []string) error {
	err := app.CSVWriter.Write(row)
	if err != nil {
		return err
	}
	app.Progress(1)
	return nil
}

func (app *App) Write(b []byte) (int, error) {
	return app.Writer.Write(b)
}

func (app *App) Print(args ...interface{}) (int, error) {
	return fmt.Fprint(app.Writer, args...)
}

func (app *App) Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(app.Writer, args...)
}

func (app *App) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(app.Writer, format, args...)
}

func (app *App) Msg(args ...interface{}) (int, error) {
	if app.Quiet {
		return 0, nil
	}
	return fmt.Fprint(os.Stderr, args...)
}

func (app *App) Msgln(args ...interface{}) (int, error) {
	if app.Quiet {
		return 0, nil
	}
	return fmt.Fprintln(os.Stderr, args...)
}

func (app *App) Msgf(format string, args ...interface{}) (int, error) {
	if app.Quiet {
		return 0, nil
	}
	return fmt.Fprintf(os.Stderr, format, args...)
}

func (app *App) Log(args ...interface{}) {
	if !app.Quiet {
		log.Print(args...)
	}
}

func (app *App) Logln(args ...interface{}) {
	if !app.Quiet {
		log.Println(args...)
	}
}

func (app *App) Logf(format string, args ...interface{}) {
	if !app.Quiet {
		log.Printf(format, args...)
	}
}
