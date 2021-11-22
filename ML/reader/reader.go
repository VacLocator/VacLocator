package reader

import (
	"encoding/csv"
	"net/http"
	m "progra_conc_TF/models"
	"strconv"

	"golang.org/x/text/encoding/charmap"
)

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(charmap.ISO8859_15.NewDecoder().Reader(resp.Body))
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getDistritoVacu(data [][]string, num_slices int) []m.CentroVacuna {
	var arrData []m.CentroVacuna
	arrDataSlices := make([][]m.CentroVacuna, num_slices)
	channels := make([]chan int, num_slices)
	for i := 0; i < num_slices; i++ {
		channels[i] = make(chan int)
		go func(channel chan int, indx int) {
			dataSlice := data[len(data)*indx/num_slices : len(data)*(indx+1)/num_slices]
			for _, row := range dataSlice {
				val0, _ := strconv.Atoi(row[0])
				val1, _ := strconv.Atoi(row[1])
				val3, _ := strconv.ParseFloat(row[3], 64)
				val4, _ := strconv.ParseFloat(row[4], 64)
				arrDataSlices[indx] = append(arrDataSlices[indx], m.CentroVacuna{
					ID_PERSONAS: val0,
					EDAD:        val1,
					NOMBRE:      row[2],
					LATITUD:     val3,
					LONGITUD:    val4,
					DISTRITO:    row[5]})
			}
			channel <- 777
		}(channels[i], i)
	}

	for _, channel := range channels {
		<-channel
	}
	for _, arrDataSlice := range arrDataSlices {
		arrData = append(arrData, arrDataSlice...)
	}
	return arrData
}

func GetDataSet(num_slices int) []m.CentroVacuna {
	url := "https://raw.githubusercontent.com/VacLocator/VacLocator/dev/Data/personas_data_aleatoria.csv"
	data, err := readCSVFromUrl(url)
	if err != nil {
		panic(err)
	}
	header_size := 7
	data = data[header_size:]
	return getDistritoVacu(data, num_slices)
}

////////////////////////////////////////////////////////////////////////////////////
func GetDatasetDistritos(data [][]string, num_slices int) []m.Distrito {
	var arrData []m.Distrito
	arrDataSlices := make([][]m.Distrito, num_slices)
	channels := make([]chan int, num_slices)
	for i := 0; i < num_slices; i++ {
		channels[i] = make(chan int)
		go func(channel chan int, indx int) {
			dataSlice := data[len(data)*indx/num_slices : len(data)*(indx+1)/num_slices]
			for _, row := range dataSlice {
				val0, _ := strconv.Atoi(row[0])
				val1, _ := strconv.Atoi(row[1])
				val3, _ := strconv.ParseFloat(row[3], 64)
				val4, _ := strconv.ParseFloat(row[4], 64)
				val6, _ := strconv.Atoi(row[6])
				arrDataSlices[indx] = append(arrDataSlices[indx], m.Distrito{
					UBIGEO:               val0,
					ID_CENTRO_VACUNACION: val1,
					NOMBRE:               row[2],
					LONGITUD:             val3,
					LATITUD:              val4,
					DISTRITO:             row[5],
					CANTIDAD:             val6})
			}
			channel <- 777
		}(channels[i], i)
	}

	for _, channel := range channels {
		<-channel
	}
	for _, arrDataSlice := range arrDataSlices {
		arrData = append(arrData, arrDataSlice...)
	}
	return arrData
}

func GetCentrosVacunaData(num_slices int) []m.Distrito {
	url := "https://raw.githubusercontent.com/VacLocator/VacLocator/dev/Data/Centros_vacuna_distritos.csv"
	data, err := readCSVFromUrl(url)
	if err != nil {
		panic(err)
	}
	header_size := 7
	data = data[header_size:]
	return GetDatasetDistritos(data, num_slices)
}
