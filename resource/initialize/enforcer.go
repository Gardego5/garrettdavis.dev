package initialize

import (
	_ "embed"

	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/jmoiron/sqlx"
)

//go:embed enforcer.model.conf
var ModelFile string

func Enforcer(db *sqlx.DB) (*casbin.Enforcer, error) {
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

	return enforcer, nil
}
