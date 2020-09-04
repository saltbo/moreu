//go:generate statik -p assets -ns moreu -src=../../moreu-front/dist -dest ..
package assets

import (
	"log"
	"net/http"

	"github.com/rakyll/statik/fs"
)

func EmbedFS() http.FileSystem {
	efs, err := fs.NewWithNamespace(Moreu)
	if err != nil {
		log.Fatalln(err)
	}

	return efs
}
