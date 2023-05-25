package models

type Client struct {
	Nombre string `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Contrasena string `json:"contrasena,omitempty" bson:"contrasena,omitempty"`
	FechaNacimiento string `json:"fecha_nacimiento,omitempty" bson:"fecha_nacimiento,omitempty"`
	Direccion string `json:"direccion,omitempty" bson:"direccion,omitempty"`
	NumeroIdentificacion string `json:"numero_identificacion,omitempty" bson:"numero_identificacion,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Telefono string `json:"telefono,omitempty" bson:"telefono,omitempty"`
	Genero string `json:"genero,omitempty" bson:"genero,omitempty"`
	Nacionalidad string `json:"nacionalidad,omitempty" bson:"nacionalidad,omitempty"`
	Ocupacion string `json:"ocupacion,omitempty" bson:"ocupacion,omitempty"`
}