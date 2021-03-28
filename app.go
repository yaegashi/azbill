package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
	"github.com/yaegashi/azbill/store"
)

const (
	ProgressScale    = 100
	ProgressColumn   = 50
	defaultClientID  = "4a034c56-da44-48ce-90db-039a408974bd"
	defaultTenantID  = "common"
	environConfigDir = "AZBILL_CONFIG_DIR"
	defaultConfigDir = "~/.azbill"
	environAuth      = "AZBILL_AUTH"
	defaultAuth      = "dev"
	environAuthFile  = "AZBILL_AUTH_FILE"
	defaultAuthFile  = "auth_file.json"
	environAuthDev   = "AZBILL_AUTH_DEV"
	defaultAuthDev   = "auth_dev.json"
	environFormat    = "AZBILL_FORMAT"
	defaultFormat    = "csv"
)

type App struct {
	Writer      io.WriteCloser
	CSVWriter   *csv.Writer
	ConfigStore *store.Store
	ConfigDir   string
	Output      string
	Auth        string
	AuthDev     string
	AuthFile    string
	Client      string
	Tenant      string
	Format      string
	Flatten     bool
	Pretty      bool
	IsStdout    bool
	Quiet       bool
	Records     int
	Column      int
	Keys        []string
	StartTime   time.Time
}

func (app *App) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "azbill",
		Short:             "Azure billing data exporter",
		PersistentPreRunE: app.PersistentPreRunE,
		SilenceUsage:      true,
		Version:           fmt.Sprintf("%s (%-0.7s)", version, commit),
	}
	cmd.PersistentFlags().StringVarP(&app.ConfigDir, "config-dir", "", "", envHelp("config dir", environConfigDir, defaultConfigDir))
	cmd.PersistentFlags().StringVarP(&app.Client, "client", "", "", envHelp("Azure client", auth.ClientID, defaultClientID))
	cmd.PersistentFlags().StringVarP(&app.Tenant, "tenant", "", "", envHelp("Azure tenant", auth.TenantID, defaultTenantID))
	cmd.PersistentFlags().StringVarP(&app.Auth, "auth", "", "", envHelp("auth source [dev,env,file,cli]", environAuth, defaultAuth))
	cmd.PersistentFlags().StringVarP(&app.AuthFile, "auth-file", "", "", envHelp("auth file store", environAuthFile, defaultAuthFile))
	cmd.PersistentFlags().StringVarP(&app.AuthDev, "auth-dev", "", "", envHelp("auth dev store", environAuthDev, defaultAuthDev))
	cmd.PersistentFlags().StringVarP(&app.Format, "format", "", "", envHelp("output format [csv,json,flatten,pretty]", environFormat, defaultFormat))
	cmd.PersistentFlags().StringVarP(&app.Output, "output", "o", "", "output file")
	cmd.PersistentFlags().BoolVarP(&app.Quiet, "quiet", "q", false, "quiet")
	return cmd
}

func envDefault(val, env, def string) string {
	if val == "" {
		val = os.Getenv(env)
	}
	if val == "" {
		val = def
	}
	return val
}

func envHelp(msg, env, def string) string {
	return fmt.Sprintf(`%s (env:%s, default:%s)`, msg, env, def)
}

func (app *App) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	app.Client = envDefault(app.Client, auth.ClientID, defaultClientID)
	app.Tenant = envDefault(app.Tenant, auth.TenantID, defaultTenantID)
	app.ConfigDir = envDefault(app.ConfigDir, environConfigDir, defaultConfigDir)
	app.Auth = envDefault(app.Auth, environAuth, defaultAuth)
	app.AuthDev = envDefault(app.AuthDev, environAuthDev, defaultAuthDev)
	app.AuthFile = envDefault(app.AuthFile, environAuthFile, defaultAuthFile)
	app.Format = envDefault(app.Format, environFormat, defaultFormat)

	store, err := store.NewStore(app.ConfigDir)
	if err != nil {
		return err
	}
	app.ConfigStore = store

	app.IsStdout = app.Output == "" || app.Output == "-"

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
			return fmt.Errorf("unknown format: %s", f)
		}
	}

	return nil
}

func (app *App) Authorize() (autorest.Authorizer, error) {
	switch app.Auth {
	case "env":
		return auth.NewAuthorizerFromEnvironment()
	case "file":
		loc, _ := app.ConfigStore.Location(app.AuthFile, true)
		app.Logf("Loading auth-file config in %s", loc)
		os.Setenv("AZURE_AUTH_LOCATION", app.AuthFile)
		return auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	case "cli":
		return auth.NewAuthorizerFromCLI()
	case "dev":
		loc, _ := app.ConfigStore.Location(app.AuthDev, true)
		app.Logf("Loading auth-dev token in %s", loc)
		b, err := app.ConfigStore.ReadFile(app.AuthDev)
		if err != nil {
			app.Logf("Warning: %s", err)
			return app.AuthorizeDeviceFlow()
		}
		var token *adal.ServicePrincipalToken
		err = json.Unmarshal(b, &token)
		if err != nil {
			app.Logf("Warning: %s", err)
			return app.AuthorizeDeviceFlow()
		}
		save := false
		token.SetRefreshCallbacks([]adal.TokenRefreshCallback{func(adal.Token) error { save = true; return nil }})
		err = token.EnsureFresh()
		if err != nil {
			app.Logf("Warning: %s", err)
			return app.AuthorizeDeviceFlow()
		}
		if save {
			b, err := json.Marshal(token)
			if err == nil {
				loc, _ := app.ConfigStore.Location(app.AuthDev, true)
				app.Logf("Saving auth-dev token in %s", loc)
				err = app.ConfigStore.WriteFile(app.AuthDev, b, 0600)
			}
			if err != nil {
				app.Logf("Warning: %s", err)
			}
		}
		return autorest.NewBearerAuthorizer(token), nil
	}
	return nil, fmt.Errorf("unknown auth: %s", app.Auth)
}

func (app *App) AuthorizeDeviceFlow() (autorest.Authorizer, error) {
	deviceConfig := auth.NewDeviceFlowConfig(app.Client, app.Tenant)
	token, err := deviceConfig.ServicePrincipalToken()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(token)
	if err == nil {
		loc, _ := app.ConfigStore.Location(app.AuthDev, true)
		app.Logf("Saving auth-dev token in %s", loc)
		err = app.ConfigStore.WriteFile(app.AuthDev, b, 0600)
	}
	if err != nil {
		app.Logf("Warning: %s", err)
	}
	return autorest.NewBearerAuthorizer(token), nil
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
		app.Writer.Write([]byte{0xef, 0xbb, 0xbf}) // UTF-8 BOM
	}
	app.Keys = nil
	app.Records = 0
	app.StartTime = time.Now()
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
	endTime := time.Now()
	d := endTime.Sub(app.StartTime)
	app.Logf("Done %d records in %s, %f records/sec", app.Records, d, float64(app.Records)/d.Seconds())
}

func (app *App) Progress(n int) {
	if app.IsStdout {
		app.Records += n
		return
	}
	if n == 0 {
		if app.Column > 0 {
			spaces := ""
			for i := app.Column; i < ProgressColumn; i++ {
				spaces += " "
			}
			app.Msgf("%s %8d records\n", spaces, app.Records)
		}
		return
	}
	records := app.Records + n
	blocks := app.Records / ProgressScale
	n = (records / ProgressScale) - blocks
	app.Records = records
	for i := 0; i < n; i++ {
		app.Msg(".")
		app.Column++
		blocks++
		if app.Column >= ProgressColumn {
			app.Msgf(" %8d records\n", blocks*ProgressScale)
			app.Column = 0
		}
	}
}

func (app *App) JSONMarshal(v interface{}) error {
	var err error
	if app.Flatten {
		v, err = flattenToMap(v, true)
		if err != nil {
			return err
		}
	}
	enc := json.NewEncoder(app.Writer)
	if app.Pretty {
		enc.SetIndent("", "  ")
	}
	err = enc.Encode(v)
	if err != nil {
		return err
	}
	app.Progress(1)
	return nil
}

func (app *App) CSVMarshal(v interface{}) error {
	if app.Keys == nil {
		m, _ := flattenToMap(v, false)
		for key := range m {
			app.Keys = append(app.Keys, key)
		}
		sort.Strings(app.Keys)
		err := app.CSVWriter.Write(app.Keys)
		if err != nil {
			return err
		}
	}
	row := []string{}
	m, err := flattenToMap(v, true)
	if err != nil {
		return err
	}
	for _, key := range app.Keys {
		col := ""
		if val, ok := m[key]; ok {
			rv := reflect.ValueOf(val)
			switch rv.Kind() {
			case reflect.Slice, reflect.Map:
				if rv.IsNil() {
					switch rv.Kind() {
					case reflect.Slice:
						col = "[]"
					default:
						col = "{}"
					}
				} else {
					b, err := json.Marshal(val)
					if err != nil {
						return err
					}
					col = string(b)
				}
			default:
				col = fmt.Sprint(val)
			}
		}
		row = append(row, col)
	}
	err = app.CSVWriter.Write(row)
	if err != nil {
		return err
	}
	app.Progress(1)
	return nil
}

func (app *App) Marshal(v interface{}) error {
	switch app.Format {
	case "csv":
		return app.CSVMarshal(v)
	default:
		return app.JSONMarshal(v)
	}
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
