package main

import (
	"fmt"
	math "math"
	m "progra_conc_TF/models"
	r "progra_conc_TF/reader"
	tsne "progra_conc_TF/tsne"

	"github.com/sjwhitworth/golearn/pca"
	"gonum.org/v1/gonum/mat"
)

type PuntosReferencia struct {
	CENTROPOB m.CentroVacuna
	DIST      float64
}

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

func parametrosEleccion(cp1 m.CentroVacuna, cp2 m.CentroVacuna) float64 {
	return math.Sqrt(math.Pow(cp1.LATITUD-cp2.LATITUD, 2) + math.Pow(cp1.LONGITUD-cp2.LONGITUD, 2))
}

func getReferencias(dataset []m.CentroVacuna, test_row m.CentroVacuna, aComparar int) []m.CentroVacuna {
	var seleccionar []PuntosReferencia
	for _, train_row := range dataset {
		dist := parametrosEleccion(test_row, train_row)
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

func PredictClassification(dataset []m.CentroVacuna, test_row m.CentroVacuna, aComparar int) string {
	references := getReferencias(dataset, test_row, aComparar)
	var output_values []string
	for _, reference := range references {
		output_values = append(output_values, reference.DISTRITO)
	}
	return mostFrequent(output_values)
}

func main() {

	num_threads_reading := 10
	arrDistritos = r.GetDataSet(num_threads_reading)

	print(PredictClassification(arrDistritos, arrDistritos[5000], 10))
}
