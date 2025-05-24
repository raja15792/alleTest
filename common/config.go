package common

type Config struct {
	DbURI           string `env:"DB_URI" envDefault:"postgres://alle:password@localhost:5432/alle_prod"`
	Port            string `env:"PORT" envDefault:"1323"`
}
