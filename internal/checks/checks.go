package checks

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/viper"
)

// Checks slice of checks with NextRun method for kicking clients to update.
type Checks []*Info

// NextRun sets the next run property of all checks to specified timestamp.
func (c Checks) NextRun(t time.Time) {
	for _, check := range c {
		check.NextRun = t
	}
}

// GetChecksFromViper gets all the checks from the viper.Viper config.
func GetChecksFromViper(cfg config.Conf, vcfg *viper.Viper, logger *slog.Logger, hostName string) Checks {
	checkListMap := make(map[string]bool)

	for _, key := range vcfg.AllKeys() {
		if strings.HasPrefix(key, "check.") {
			token := strings.SplitN(key, ".", 3) //nolint:mnd // check keys come in 3 parts.
			checkListMap[token[1]] = true
		}
	}

	checkList := make(Checks, len(checkListMap))
	i := 0

	for v := range checkListMap {
		checkList[i] = GetCheckFromViper(cfg, logger, v, hostName)
		logger.Info("adding check", slog.String("check.name", checkList[i].Name),
			slog.Duration("check.period", checkList[i].Period))
		i++
	}

	return checkList
}

// GetCheckFromViper returns a check with the specified name from the config file.
func GetCheckFromViper(cfg config.Conf, logger *slog.Logger, name, hostName string) *Info {
	check := &Info{
		Name:     cfg.GetString(fmt.Sprintf("check.%s.name", name)),
		Period:   cfg.GetDuration(fmt.Sprintf("check.%s.period", name)),
		Command:  cfg.GetString(fmt.Sprintf("check.%s.command", name)),
		Hostname: hostName,
		Timeout:  cfg.GetDuration(fmt.Sprintf("check.%s.timeout", name)),
		Workdir:  cfg.GetString(fmt.Sprintf("check.%s.workdir", name)),
	}

	if check.Timeout == 0 {
		check.Timeout = cfg.GetDuration("default.timeout")
	}

	if check.Period == 0 {
		check.Period = cfg.GetDuration("default.period")
	}

	if check.Name == "" {
		switch strings.ToLower(cfg.GetString("default.name-format")) {
		case "lowercase", "lower", "lc":
			check.Name = strings.ToUpper(name)
		case "titlecase", "title", "tc":
			check.Name = strings.ToTitle(name)
		case "uppercase", "upper", "uc":
			fallthrough
		default:
			check.Name = strings.ToUpper(name)
		}
	}

	switch cfg.GetString(fmt.Sprintf("check.%s.type", name)) {
	case "host":
		check.Type = api.CheckType_HOST
	case "", "service":
		check.Type = api.CheckType_SERVICE
	default:
		logger.Warn("unknown check type, defaulting to 'service'",
			slog.String("check", check.Name),
			slog.String("check-type-supplied", cfg.GetString(fmt.Sprintf("check.%s.type", name))),
		)

		check.Type = api.CheckType_SERVICE
	}

	return check
}
