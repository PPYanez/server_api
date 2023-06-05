package menu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"trust-bank/api/db"
	"trust-bank/api/models"

	"go.mongodb.org/mongo-driver/bson"
)

func BankMenu(nroIdentificacion string) {
	var keepRunning = true

	for keepRunning {
		fmt.Print(
			"1. Realizar deposito\n",
			"2. Realizar transferencia\n",
			"3. Realizar giro \n",
			"4. Salir\n",
			"Ingrese una opción: ",
		)

		var option int
		fmt.Scanln(&option)

		if option == 1 {
			mongoClient := db.GetClient()

			walletColl := mongoClient.Database("trustBank").Collection("Billeteras")
			var walletResult models.Wallet
			err := walletColl.FindOne(context.TODO(), bson.M{"nro_cliente": nroIdentificacion}).Decode(&walletResult)

			if err != nil {
				fmt.Println("La cuenta ingresada no tiene billetera asociada")
				continue
			}
			var monto string

			fmt.Print("Ingrese el monto: ")
			fmt.Scanln(&monto)

			depositInfo := models.Deposit{
				NumeroCliente: nroIdentificacion,
				Monto:         monto,
				Divisa:        "USD",
			}

			depositInfoJson, _ := json.Marshal(depositInfo)
			url, err := url.Parse("http://localhost:8080/api/deposito")
			if err != nil {
				panic("URL invalida")
			}

			req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(depositInfoJson))

			if err != nil {
				panic("Error al crear petición")
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()

			if err != nil {
				panic("Error al enviar petición")
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic("Error al leer respuesta")
			}

			if string(body) == `{"estado":"deposito_enviado"}` {
				fmt.Println("El deposito ha sido enviado correctamente")
			}

			if string(body) == `{"estado":"cliente_no_encontrado"}` {
				fmt.Println("El cliente no fue encontrado")
			}

			if string(body) == `{"estado":"billetera_no_encontrada"}` {
				fmt.Println("La billetera no fue encontrada")
			}
		}

		if option == 2 {
			mongoClient := db.GetClient()

			walletColl := mongoClient.Database("trustBank").Collection("Billeteras")
			var walletOriginResult models.Wallet
			err := walletColl.FindOne(context.TODO(), bson.M{"nro_cliente": nroIdentificacion}).Decode(&walletOriginResult)

			if err != nil {
				fmt.Println("La cuenta ingresada no tiene billetera asociada")
				continue
			}
			var nro_cliente_destino string

			fmt.Print("Ingrese el número de cliente destino: ")
			fmt.Scanln(&nro_cliente_destino)

			coll := mongoClient.Database("trustBank").Collection("Clientes")

			var result models.Client
			err = coll.FindOne(context.TODO(), bson.M{"numero_identificacion": nro_cliente_destino}).Decode(&result)

			if err != nil {
				fmt.Println("La cuenta ingresada no existe")
				continue
			}

			var walletResult models.Wallet
			err = walletColl.FindOne(context.TODO(), bson.M{"nro_cliente": nro_cliente_destino}).Decode(&walletResult)

			if err != nil {
				fmt.Println("La cuenta ingresada no tiene billetera asociada")
				continue
			}

			var monto string

			fmt.Print("Ingrese el monto: ")
			fmt.Scanln(&monto)

			transferInfo := models.Transfer{
				NumeroClienteOrigen:  nroIdentificacion,
				NumeroClienteDestino: nro_cliente_destino,
				Monto:                monto,
				Divisa:               "USD",
			}

			transferInfoJson, _ := json.Marshal(transferInfo)
			url, err := url.Parse("http://localhost:8080/api/transferencia")
			if err != nil {
				panic("URL invalida")
			}

			req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(transferInfoJson))

			if err != nil {
				panic("Error al crear petición")
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()

			if err != nil {
				panic("Error al enviar petición")
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic("Error al leer respuesta")
			}

			if string(body) == `{"estado":"transferencia_enviada"}` {
				fmt.Println("La transferencia ha sido enviada correctamente")
			}

			if string(body) == `{"estado":"billetera_origen_sin_fondos_suficientes"}` {
				fmt.Println("Su saldo es insuficiente")
			}
		}

		if option == 3 {
			mongoClient := db.GetClient()

			walletColl := mongoClient.Database("trustBank").Collection("Billeteras")
			var walletResult models.Wallet
			err := walletColl.FindOne(context.TODO(), bson.M{"nro_cliente": nroIdentificacion}).Decode(&walletResult)

			if err != nil {
				fmt.Println("La cuenta ingresada no tiene billetera asociada")
				continue
			}

			var monto string

			fmt.Print("Ingrese un monto: ")
			fmt.Scanln(&monto)

			withdrawInfo := models.Deposit{
				NumeroCliente: nroIdentificacion,
				Monto:         monto,
				Divisa:        "USD",
			}

			withdrawInfoJson, _ := json.Marshal(withdrawInfo)
			url, err := url.Parse("http://localhost:8080/api/giro")
			if err != nil {
				panic("URL invalida")
			}

			req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(withdrawInfoJson))

			if err != nil {
				panic("Error al crear petición")
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()

			if err != nil {
				panic("Error al enviar petición")
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic("Error al leer respuesta")
			}

			if string(body) == `{"estado":"giro_enviado"}` {
				fmt.Println("El giro ha sido solicitado correctamente")
			}

			if string(body) == `{"estado":"cliente_no_encontrado"}` {
				fmt.Println("El cliente no fue encontrado")
			}

			if string(body) == `{"estado":"billetera_origen_sin_fondos_suficientes"}` {
				fmt.Println("Su saldo es insuficiente")
			}

		}

		if option == 4 {
			keepRunning = false
		}
	}
}
