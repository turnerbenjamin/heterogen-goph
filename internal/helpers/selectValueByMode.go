package helpers

import "log"

type mode string

const Development = mode("development")
const Production = mode("production")

type ValueSelector[t any] struct {
	Development t
	Production  t
}

func NewMode(m string) mode {
	switch m {
	case "development":
		{
			return Development
		}
	case "production":
		{
			return Production
		}
	default:
		{
			log.Fatal("Invalid mode selected")
			return Production
		}
	}
}

func SelectValueByMode[t any](m mode, vs ValueSelector[t]) t {
	switch m {
	case Development:
		{
			return vs.Development
		}
	case Production:
		{
			return vs.Production
		}
	default:
		{
			log.Fatal("Invalid mode selected")
			return vs.Production
		}
	}
}
