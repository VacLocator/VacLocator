package main

import (
	"fmt"
	"io"
	math "math"
	"os"
	m "progra_conc_TF/models"
	r "progra_conc_TF/reader"
	tsne "progra_conc_TF/tsne"

	//"github.com/kniren/gota/dataframe"
	"github.com/sjwhitworth/golearn/pca"
	"gonum.org/v1/gonum/mat"

	"net/http"
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

/*
func requestearDatos(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	lat, ok1 := vars["lat"]

	lng, ok2 := vars["lng"]

	cod1, err1 := strconv.Atoi(lat)

	cod2, err2 := strconv.Atoi(lng)

	fmt.Print(cod1 + cod2)
	if ok2 == nil {
		fmt.Print(ok2)
	}
	if err1 != nil {
		fmt.Print(err1)
	}
	if ok1 == nil {
		fmt.Print(ok1)
	}
	if err2 != nil {
		fmt.Print(err2)
	} else {
		log.Println(cod1)
		response.Header().Set("Content-Type", "application/json	")

		var oTraceGo response

		jsonBytes, _ := json.MarshalIndent(oTraceGo, "", " ")

		log.Println(string(jsonBytes))

		io.WriteString(response, string(jsonBytes))
	}

}
*/

func main() {
	/*
		url :="file:///Z:/TFConcu/gitjab/VacLocator/Frontend/coord.html?lat=&lng="
		params := (new url(url)).searchParams
		latitud:=params.get('lat') // "n1"
		longitud:=params.get('lng')
	*/

	//--
	response, err := http.Get("file:///Z:/TFConcu/gitjab/VacLocator/Frontend/coord.html?lat=&lng=") //use package "net/http"

	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	// Copy data from the response to standard output
	n, err1 := io.Copy(os.Stdout, response.Body) //use package "io" and "os"
	if err != nil {
		fmt.Println(err1)
		return
	}

	fmt.Print(n)

	//--

	ls := []string{}
	n_cantidades := []int{}
	cantidad := true
	num_threads_reading := 10
	arrDistritos = r.GetDataSet(num_threads_reading)
	//----
	/*
		latitud = requestearDatos(cantidad, r)
		longitud = requestearDatos(cantidad, r)
	*/

	/*
		latitud_nueva := -12.06702829
		longitud_nueva := -77.0114123
	*/
	//----
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

	fmt.Printf("%v", ls)

}
