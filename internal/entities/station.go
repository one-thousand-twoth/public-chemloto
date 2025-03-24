package entities

import "github.com/anrew1002/Tournament-ChemLoto/internal/common"

type InitFunction func(chan common.Message)

type Channel struct {
	ID   ID                        `json:"-"`
	Name string                    `json:"-"`
	Type string                    `json:"-"`
	Fn   func(chan common.Message) `json:"-"`
}
