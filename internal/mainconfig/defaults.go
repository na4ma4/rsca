package mainconfig

import "github.com/spf13/viper"

// ConfigInit is the common config initialisation for the commands.
//
//nolint:mnd // defaults are magic.
func ConfigInit() {
	viper.SetConfigName("rsca")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./artifacts")
	viper.AddConfigPath("./testdata")
	viper.AddConfigPath("$HOME/.rsca")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath("/etc/rsca")
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath("/usr/local/rsca/etc")
	viper.AddConfigPath("/opt/homebrew/etc")
	viper.AddConfigPath("/etc/nsca")
	viper.AddConfigPath("/etc/nagios")
	viper.AddConfigPath(".")

	viper.SetDefault("general.jitter", "10s")
	viper.SetDefault("general.retry", true)
	viper.SetDefault("general.max-retries", 3)
	viper.SetDefault("general.check-tick", "9s")
	viper.SetDefault("general.tags", []string{})
	viper.SetDefault("general.registration-interval", "180s")

	viper.SetDefault("default.period", "120s")
	viper.SetDefault("default.timeout", "3s")
	viper.SetDefault("default.name-format", "uppercase")

	viper.SetDefault("nagios.command-file", "/tmp/nagios.cmd")

	viper.SetDefault("admin.server", "127.0.0.1:15888")
	viper.SetDefault("admin.cert-type", "Cert")

	viper.SetDefault("client.server", "127.0.0.1:15888")
	viper.SetDefault("client.cert-type", "Client")

	viper.SetDefault("server.listen", "0.0.0.0:15888")
	viper.SetDefault("server.tick", "15s")
	viper.SetDefault("server.cert-type", "Server")
	viper.SetDefault("server.state-store", "/tmp/rsca-state.db")
	viper.SetDefault("server.state-timeout", "120s")
	viper.SetDefault("server.state-tick", "60s")

	viper.SetDefault("watchdog.enabled", false)
	viper.SetDefault("watchdog.tick", "30s")

	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.listen", "localhost:2112")
	viper.SetDefault("metrics.timeout.read", "1h")
	viper.SetDefault("metrics.timeout.read-header", "10s")
	viper.SetDefault("metrics.timeout.write", "1m")
	viper.SetDefault("metrics.timeout.idle", "10s")

	_ = viper.ReadInConfig()
}
