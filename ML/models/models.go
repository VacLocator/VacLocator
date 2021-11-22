package models

type CentroVacuna struct {
	//Id del centro de vacunación
	ID_PERSONAS int `json:"ubigeo"`
	//Nombre del centro de vacunación
	EDAD int `json:"edad"`
	//Id del centro de vacunación
	NOMBRE string `json:"nombre"`
	// Coordenada: Latitud;
	LONGITUD float64 `json:"longitud"`
	// Coordenada: Longitud;
	LATITUD float64 `json:"latitud"`
	// Distrito de procedencia
	DISTRITO string `json:"distrito"`
}

type Prediction struct {
	CODCP  int    `json:"codcp"`
	EXPECT string `json:"expect"`
	PRED   string `json:"pred"`
}
