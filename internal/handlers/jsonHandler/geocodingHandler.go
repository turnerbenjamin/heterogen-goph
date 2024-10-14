package jsonHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
	"github.com/turnerbenjamin/heterogen-go/internal/models"
)

type geoCodingResponse struct {
	ErrorMessage string `json:"error_message"`
	Results      []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

var GeoCodingHandler = func(w http.ResponseWriter, r *http.Request, m *models.ResponseModel) error {
	var token = os.Getenv("GOOGLE_GEOCODING_TOKEN")
	postcode := r.URL.Query().Get("postcode")
	postcode = strings.ReplaceAll(postcode, " ", "")
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?components=postal_code:%s&key=%s", postcode, token)
	//Make API Request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	//Parse Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Bad request: " + url)
	}

	respStruct := geoCodingResponse{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return err
	}

	if respStruct.ErrorMessage != "" {
		return errors.New(respStruct.ErrorMessage)
	}

	if len(respStruct.Results) == 0 {
		return httpErrors.Make(http.StatusNotFound, []httpErrors.ErrorMessage{"Postcode not found"})
	}

	//Return Lat and Lng
	location, err := json.Marshal(respStruct.Results[0].Geometry.Location)
	if err != nil {
		return httpErrors.ServerFail()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(location)

	return nil
}
