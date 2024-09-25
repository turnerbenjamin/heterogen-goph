package hg_services

import (
	"database/sql"
	"strings"

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
	(id, reference, trading_name, location, postcode, is_grower, cph_number, logo, about, email_address, website)
	VALUES($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, reference, trading_name, location, postcode, is_grower, cph_number, logo, about, email_address, website
	;
	`
	var nb models.Business
	rows, err := bs.db.Query(baseQuery, bns.Id, bns.Reference, bns.TradingName, bns.Location, bns.Postcode, bns.IsGrower, bns.CPH_Number, bns.Logo, bns.About, bns.EmailAddress, bns.Website)

	if err != nil {
		if strings.Contains(err.Error(), "businesses_trading_name_key") {
			return nil, httpErrors.Make(401, []httpErrors.ErrorMessage{"Trading name is already associated with a business"})
		}
		return nil, httpErrors.ServerFail()
	}

	if rows.Next() {
		err := rows.Scan(&nb.Id, &nb.Reference, nb.TradingName, nb.Location, nb.Postcode, nb.IsGrower, nb.CPH_Number, bns.Logo, bns.About, bns.EmailAddress, bns.Website)
		if err != nil {
			return nil, httpErrors.ServerFail()
		}
	}
	return &nb, nil
}
