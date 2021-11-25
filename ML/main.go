package main

import (
	"fmt"
	"io"
	"log"
	math "math"
	"net/http"
	m "progra_conc_TF/models"
	r "progra_conc_TF/reader"
	tsne "progra_conc_TF/tsne"

	//"github.com/kniren/gota/dataframe"
	"encoding/json"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sjwhitworth/golearn/pca"
	"gonum.org/v1/gonum/mat"
)

type PuntosReferencia struct {
	CENTROPOB m.CentroVacuna
	DIST      float64
}

var latitud float64

var longitud float64

var arrDistritos []m.CentroVacuna

func mostFrequent(arr []string) string {
	m := map[string]int{}
	var maxCnt int
	var freq string
	for _, a := range arr {
		m[a]++
		if m[a] > maxCnt {
			maxCnt = m[a]
			freq = a
		}
	}
	return freq
}

func parametrosEleccion(lat float64, lon float64, cp2 m.CentroVacuna) float64 {
	return math.Sqrt(math.Pow(lat-cp2.LATITUD, 2) + math.Pow(lon-cp2.LONGITUD, 2))
}

func getReferencias(dataset []m.CentroVacuna, lat float64, lon float64, aComparar int) []m.CentroVacuna {
	var seleccionar []PuntosReferencia
	for _, train_row := range dataset {
		dist := parametrosEleccion(lat, lon, train_row)
		seleccionar = append(seleccionar, PuntosReferencia{CENTROPOB: train_row, DIST: dist})
	}
	var neighbors []m.CentroVacuna
	neighbor_indxes := make([]int, aComparar)
	var min float64 = 9999999
	for indx, distance := range seleccionar {
		if distance.DIST < min {
			min = distance.DIST
			neighbor_indxes = append(neighbor_indxes, indx)
			if len(neighbor_indxes) > aComparar {
				neighbor_indxes = neighbor_indxes[1:]
			}
		}
	}

	X := mat.NewDense(3, 5, nil)

	perplexity := float64(300)
	learningRate := float64(300)
	pcaComponents := 50

	Xdense := mat.DenseCopyOf(X)
	pcaTransform := pca.NewPCA(pcaComponents)
	Xt := pcaTransform.FitTransform(Xdense)

	t := tsne.NewTSNE(2, perplexity, learningRate, 300, true)
	t.EmbedData(Xt, func(iter int, divergence float64, embedding mat.Matrix) bool {
		if iter%10 == 0 {
			fmt.Printf("Iteration %d: divergence is %v\n", iter, divergence)
		}
		return false
	})

	for _, neighbor_indx := range neighbor_indxes {
		neighbors = append(neighbors, seleccionar[neighbor_indx].CENTROPOB)
	}
	return neighbors
}

func PredictClassification(dataset []m.CentroVacuna, lat float64, lon float64, aComparar int) string {
	references := getReferencias(dataset, lat, lon, aComparar)
	var output_values []string
	for _, reference := range references {
		output_values = append(output_values, reference.DISTRITO)
	}
	return mostFrequent(output_values)
}

func manejadorSolicitudes() {
	//enrutador
	r := mux.NewRouter()
	enableCORS(r)
	//endpoints
	r.HandleFunc("/lat/{latitud}/lng/{longitud}/filtro/{filtro}", buscarCercano)
	//r.HandleFunc("/agregar", agregarCentroDeVacunacion)

	log.Fatal(http.ListenAndServe(":9000", r))
}

func buscarCercano(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	lat, ok := vars["latitud"]
	lng, ok := vars["longitud"]
	filtro, ok := vars["filtro"]

	fmt.Print("filtro : ")
	fmt.Println(filtro)

	if latitud, err := strconv.ParseFloat(lat, 64); err == nil {
		if longitud, err := strconv.ParseFloat(lng, 64); err == nil {
			if ok {

				fmt.Print("latitud : ")
				fmt.Println(latitud)

				fmt.Print("longitud : ")
				fmt.Println(longitud)

				ls := []string{}
				n_cantidades := []int{}
				cantidad := false
				if filtro == "1" {
					cantidad = true
				}
				num_threads_reading := 10
				arrDistritos = r.GetDataSet(num_threads_reading)
				distrito_previo := PredictClassification(arrDistritos, latitud, longitud, 10)

				centros_vacunacion := r.GetCentrosVacunaData(num_threads_reading)
				for _, row := range centros_vacunacion {
					if distrito_previo == row.DISTRITO {
						ls = append(ls, row.NOMBRE)
						n_cantidades = append(n_cantidades, row.CANTIDAD)
					}
				}
				if cantidad {
					n := len(ls)
					if n > 1 {
						swapped := true
						for swapped {
							swapped = false

							for i := 0; i < n-1; i++ {

								if n_cantidades[i] > n_cantidades[i+1] {

									n_cantidades[i], n_cantidades[i+1] = n_cantidades[i+1], n_cantidades[i]
									ls[i], ls[i+1] = ls[i+1], ls[i]

									swapped = true
								}
							}
						}
					}
				}

				//To DO: Envia BC lista
				//To Do: lista = Respuesta
				//To Do: Recibe la lista de BC

				centros := r.GetCentrosVacunaData(num_threads_reading)

				fmt.Print("centros : ")
				fmt.Println(centros)

				var centrosMasCercano []m.Distrito

				for _, current := range ls {
					for _, row := range centros {
						if current == row.NOMBRE {
							centrosMasCercano = append(centrosMasCercano, row)
						}
					}
				}

				fmt.Print("distritoMasCercano : ")
				fmt.Println(centrosMasCercano)

				jsonBytes, _ := json.MarshalIndent(centrosMasCercano, "", " ")
				io.WriteString(response, string(jsonBytes))
			}
		}
	}

}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			// Just put some headers to allow CORS...
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			// and call next handler!
			next.ServeHTTP(w, req)
		})
}

func main() {
	manejadorSolicitudes()
}
