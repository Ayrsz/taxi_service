module github.com/gjcms/taxi_service

go 1.23.5 // Certifique-se de que é a versão correta do seu Go

require (
	github.com/cucumber/godog v0.15.0
	// Mantenha todas as suas dependências externas aqui, como antes
	github.com/gofiber/fiber/v2 v2.52.9
	github.com/golang-jwt/jwt/v5 v5.2.3
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.17.0
	gorm.io/driver/postgres v1.5.11
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// ATUALIZE ESTAS DIRETIVAS 'replace' para refletir suas NOVAS pastas na raiz
// E remova as antigas que apontavam para 'pkg/...'
replace github.com/gjcms/taxi_service/utils => ./utils

replace github.com/gjcms/taxi_service/config => ./config

replace github.com/gjcms/taxi_service/controllers => ./controllers // Antigo 'handlers'

replace github.com/gjcms/taxi_service/database => ./database

replace github.com/gjcms/taxi_service/models => ./models // AGORA 'models' está na raiz

replace github.com/gjcms/taxi_service/middlewares => ./middlewares // Antigo 'middleware'

replace github.com/gjcms/taxi_service/routes => ./routes

replace github.com/gjcms/taxi_service/services => ./services // Antigo 'utils'

// Certifique-se de que NÃO há mais `replace` apontando para `pkg/...`
// Certifique-se de que NÃO há `replace` para `/database/models` se `models` foi movido para a raiz.
