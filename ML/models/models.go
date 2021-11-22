package models

type CentroVacuna struct {
	//Id del centro de vacunación
	ID_PERSONAS int `json:"id_persona"`
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
	//Id del centro de vacunación

}

type Distrito struct {
	//Id del centro de vacunación
	UBIGEO int `json:"ubigeo"`
	//Id del centro de vacunación
	ID_CENTRO_VACUNACION int `json:"id_centro"`
	//Nombre del centro de vacunación
	NOMBRE string `json:"nombre_centro"`
	// Coordenada: Latitud;
	LONGITUD float64 `json:"longitud"`
	// Coordenada: Longitud;
	LATITUD float64 `json:"latitud"`
	// Distrito de procedencia
	DISTRITO string `json:"distrito"`
	// cantidad aceptada por centro
	CANTIDAD int `json:"cantidad"`
}

type Prediction struct {
	CODCP  int    `json:"codcp"`
	EXPECT string `json:"expect"`
	PRED   string `json:"pred"`
}
