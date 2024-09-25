package jsonHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?components=postal_code:%s&key=%s", postcode, token)

	//Make API Request
	resp, err := http.Get(url)
	if err != nil {
		return httpErrors.ServerFail()
	}

	//Parse Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return httpErrors.ServerFail()
	}

	respStruct := geoCodingResponse{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return httpErrors.ServerFail()
	}

	if respStruct.ErrorMessage != "" {
		return httpErrors.ServerFail()
	}

	if len(respStruct.Results) == 0 {
		return httpErrors.Make(http.StatusNotFound, []httpErrors.ErrorMessage{"Postcode not found"})
	}

	//Return Lat and Lng
	location, err := json.Marshal(respStruct.Results[0].Geometry.Location)
	if err != nil {
		return httpErrors.ServerFail()
	}

	w.Header().Set("content-Type", "application/json")
	w.Write(location)
	return nil
}
