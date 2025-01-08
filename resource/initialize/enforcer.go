package initialize

import (
	_ "embed"
	"log/slog"

	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/jmoiron/sqlx"
)

//go:embed enforcer.model.conf
var ModelFile string

func Enforcer(db *sqlx.DB, slogLogger *slog.Logger) (*casbin.Enforcer, error) {
	adapter, err := sqlxadapter.NewAdapter(db, "")
	if err != nil {
		return nil, err
	}

	model, err := model.NewModelFromString(ModelFile)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(model, adapter)
	if err != nil {
		return nil, err
	}

	enforcer.EnableLog(true)

	enforcer.SetLogger(&logger{logger: slogLogger})

	return enforcer, nil
}

type logger struct {
	enabled bool
	logger  *slog.Logger
}

func (l *logger) EnableLog(enable bool) { l.enabled = enable }
func (l *logger) IsEnabled() bool       { return l.enabled }
func (l *logger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	l.logger.Info("Enforce", "matcher", matcher, "request", request, "result", result, "explains", explains)
}
func (l *logger) LogModel(model [][]string) {
	l.logger.Info("Model", "model", model)
}

func (l *logger) LogRole(roles []string) {
	l.logger.Info("Role", "roles", roles)
}

func (l *logger) LogPolicy(policy map[string][][]string) {
	l.logger.Info("Policy", "policy", policy)
}
func (l *logger) LogError(err error, msg ...string) {
	l.logger.Error("Error", "error", err, "msg", msg)
}
