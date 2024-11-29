//go:generate msgp
package model

type Subject struct {
	User string `msg:"user"`
}
