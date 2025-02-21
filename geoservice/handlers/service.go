package handlers

import (
	"encoding/json"
	"fmt"
	"geoservice/metrics"
	"log"
	"net/http"
)

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}

func Authorization() *GeoService {
	ApiKeyValue := "37d566ad6b5aea851c57d0598374b069da969402"
	SecretKeyValue := "38c8b602612cf1cd0e827b247e11651876601e41"
	return NewGeoService(ApiKeyValue, SecretKeyValue)
}

// @Summary Поиск адреса
// @Description Поиск адреса по текстовому запросу
// @Tags addresses
// @Accept json
// @Produce json
// @Param request body SearchRequest true "Поисковый запрос"
// @Success 200 {object} SearchResponse
// @Router /address/search [post]
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Замеряем время выполнения эндпоинта
	defer metrics.MeasureEndpointDuration("/api/address/search")()

	// Замеряем время обращения к внешнему API
	defer metrics.MeasureExternalAPIAccess("search")()

	log.Println("SearchHandler called")

	var request SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	geoService := Authorization()
	log.Printf("geoService authorization data: %v\n", geoService)

	addresses, err := geoService.AddressSearch(request.Query)
	if err != nil {
		log.Printf("Address search error: %v", err)
		http.Error(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	response := SearchResponse{Addresses: addresses}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("error encode: ", err)
	}
}

// @Summary Геокодирование адреса
// @Description Получение координат по адресу
// @Tags addresses
// @Accept json
// @Produce json
// @Param request body SearchRequest true "Адрес"
// @Success 200 {object} SearchResponse
// @Router /address/geocode [post]
func GeoHandler(w http.ResponseWriter, r *http.Request) {
	// Замеряем время выполнения эндпоинта
	defer metrics.MeasureEndpointDuration("/api/address/geocode")()

	// Замеряем время обращения к внешнему API
	defer metrics.MeasureExternalAPIAccess("geocode")()

	log.Println("GeoHandler called")

	var request GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	log.Printf("request geocode:%v\n", request)

	geoService := Authorization()
	log.Printf("geoService authorization data: %v\n", geoService)

	addresses, err := geoService.GeoCode(request.Lat, request.Lng)
	if err != nil {
		log.Printf("Geocode error: %v", err)
		http.Error(w, "Failed to fetch geocoded addresses", http.StatusInternalServerError)
		return
	}

	response := GeocodeResponse{Addresses: addresses}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("error encode: ", err)
	}

	for _, address := range response.Addresses {
		fmt.Printf("Выбран адрес: City: %s, Street: %s, House: %s, Lat: %s, Lon: %s\n",
			address.City, address.Street, address.House, address.Lat, address.Lon)
	}
}
