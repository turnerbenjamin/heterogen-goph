package hg_services

import (
	"database/sql"
	"strings"

	"github.com/lib/pq"
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
)

type BusinessService interface {
	Create(models.Business) (*models.Business, error)
}

type businessService struct {
	db *sql.DB
}

func NewBusinessServiceService(database *sql.DB) BusinessService {
	return &businessService{
		db: database,
	}
}

/*
CREATE A NEW BUSINESS
*/
func (bs *businessService) Create(bns models.Business) (*models.Business, error) {
	baseQuery := `
	INSERT INTO businesses
	(id, reference, trading_name, location, address, postcode, is_grower, cph_number, email_address)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, reference, trading_name, location, address, postcode, is_grower, cph_number, email_address
	;
	`
	var nb models.Business
	rows, err := bs.db.Query(baseQuery, bns.Id, bns.Reference, bns.TradingName, bns.Location, pq.Array(bns.Address), bns.Postcode, bns.IsGrower, bns.CphNumber, bns.EmailAddress)

	if err != nil {
		if strings.Contains(err.Error(), "businesses_trading_name_key") {
			return nil, httpErrors.Make(401, []httpErrors.ErrorMessage{"Trading name is already associated with a business"})
		}
		return nil, err
	}

	if rows.Next() {
		err := rows.Scan(&nb.Id, &nb.Reference, nb.TradingName, nb.Location, nb.Address, nb.Postcode, nb.IsGrower, nb.CphNumber, bns.EmailAddress)
		if err != nil {
			return nil, err
		}
	}
	return &nb, nil
}
