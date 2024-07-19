package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse interface {
	LocationArea |
	ExploreArea |
	Pokemon
}

func PaseJSON[T JsonResponse](jsonByte []byte) (T, bool) {
	var data T

	jsonErr := json.Unmarshal(jsonByte, &data)

	if jsonErr != nil {
		return data, false
	}

	return data, true
}

func MakeGetCall(url string) ([]byte, error) {
	res, err := http.Get(url)

	if err != nil {
		return []byte(""), errors.New("something went wrong with the network")
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return []byte(""), errors.New("something went wrong while passing data")
	}

	return body, nil
}