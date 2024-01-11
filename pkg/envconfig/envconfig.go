package envconfig

import (
	"os"
	"reflect"
	"strings"
)

// GetMissingEnvVars:
// Recibe una estructura de configuración y devuelve una lista de variables de entorno faltantes.
// La función utiliza reflexión para recorrer la estructura de configuración y verificar si los
// campos etiquetados como "required" tienen valores correspondientes en las variables de entorno.
// Si se encuentra una estructura anidada, la función se llama recursivamente para checar los
// campos de la estructura anidada. La función devuelve una lista de cadenas que contiene los
// nombres de las variables de entorno faltantes.
func GetMissingEnvVars(configuration any) []string {
	value := reflect.ValueOf(configuration).Elem()
	var missingVars []string
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			if mev := GetMissingEnvVars(field.Addr().Interface()); len(mev) > 0 {
				missingVars = append(missingVars, mev...)
			}
		default:
			envVarTag := value.Type().Field(i).Tag.Get("env")
			if envVarTag == "" || !strings.Contains(envVarTag, "required") {
				continue
			}
			key := strings.Split(envVarTag, ",")[0]
			if envVarValue, exists := os.LookupEnv(key); !exists || envVarValue == "" {
				missingVars = append(missingVars, key)
			}
		}
	}
	return missingVars
}
