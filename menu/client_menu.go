package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	menu "trust-bank/api/menu/bankMenu"
	"trust-bank/api/models"
)

func main() {
	keepRunning := true

	fmt.Println("Bienvenido a TrustBank!")

	for keepRunning {
		fmt.Print(
			"1. Iniciar sesión\n",
			"2. Salir\n",
			"Ingrese una opción: ",
		)

		var option int
		fmt.Scanln(&option)

		if option == 1 {

			var id string
			var contrasena string

			fmt.Print("Ingrese su número de identificación: ")
			fmt.Scanln(&id)
			fmt.Print("Ingrese su contraseña: ")
			fmt.Scanln(&contrasena)

			//verify login
			//if login is correct, show menu
			//if login is incorrect, show error message

			loginInfo := models.Client{
				NumeroIdentificacion: id,
				Contrasena:           contrasena,
			}

			loginInfoJson, _ := json.Marshal(loginInfo)

			url, err := url.Parse("http://localhost:8080/api/inicio_sesion")
			if err != nil {
				panic("URL invalida")
			}

			req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(loginInfoJson))

			if err != nil {
				panic("Error al crear petición")
			}

			client := http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()

			if err != nil {
				panic("Error al hacer petición")
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic("Error al leer respuesta")
			}

			if string(body) == `{"estado":"exitoso"}` {
				fmt.Println("Login exitoso!")
				menu.BankMenu(id)

			}
			if string(body) == `{"estado":"no_exitoso"}` {
				fmt.Println("Número de identificación o contraseña incorrecta")
			}
		}
		if option == 2 {
			keepRunning = false
			fmt.Println("Gracias por usar TrustBank!")
		}
	}
}
