package clients

import (
	"encoding/json"
	"errors"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/clients/responses"
	"io/ioutil"
	"log"
	"net/http"
)

type AuthClientInterface interface {
	GetUserInfo(token string) (*responses.UserInfo, error)
}

type AuthClient struct {
}

func NewAuthClient() *AuthClient {
	return &AuthClient{}
}

func (auth *AuthClient) GetUserInfo(token string) (*responses.UserInfo, error) {
	apiUrl := "http://w230847.ferozo.com/tp_prog2/api/Account/UserInfo"
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Printf("Error al crear la solicitud GET: %v", err)
		return nil, err
	}

	req.Header.Add("Authorization", token)
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error al realizar la solicitud GET: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error al leer el cuerpo de la respuesta: %v", err)
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Printf("Error al realizar la solicitud GET, código de estado: %d, cuerpo de la respuesta: %s", response.StatusCode, responseBody)
		return nil, errors.New("La peticion respondio con error")
	}

	bodyString := string(responseBody)
	log.Printf("Cuerpo de la respuesta: %s", bodyString)

	var userInfo responses.UserInfo
	if err := json.Unmarshal([]byte(bodyString), &userInfo); err != nil {
		log.Printf("Error al deserializar el JSON: %v", err)
		return nil, err
	}

	userInfo.UserId = userInfo.Codigo

	log.Printf("Código de estado: %s", response.Status)
	log.Printf("Información del usuario obtenida: %+v", userInfo)

	return &userInfo, nil
}
