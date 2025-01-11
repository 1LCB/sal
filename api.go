package sal

import (
	"log"
	"net/http"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-http/swagger"
)

var (
	Swag *swagno.Swagger
)

type API struct {
	Title string
}

func NewAPI(title string) *API {
	Swag = swagno.New(swagno.Config{
		Title:   title,
		Version: "1.0.0",
	})

	return &API{
		Title: title,
	}
}

func (x *API) Run(host string) {
	http.HandleFunc("/swagger/", swagger.SwaggerHandler(Swag.MustToJson()))

	log.Printf("| %s Running on %s\n", x.Title, host)
	log.Println("| Swagger docs available on /swagger/index.html")

	if err := http.ListenAndServe(host, nil); err != nil {
		panic(err)
	}
}
