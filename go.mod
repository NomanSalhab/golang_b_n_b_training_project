module github.com/NomanSalhab/golang_b_n_b_training_project

go 1.19

replace github.com/nomansalhab/golang_b_n_b_training_project/cmd/web/ => ./cmd/web/

require (
	github.com/alexedwards/scs/v2 v2.5.0
	github.com/go-chi/chi v1.5.4
	github.com/justinas/nosurf v1.1.1
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/jackc/pgconn v1.13.0
	github.com/jackc/pgx/v5 v5.0.3
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/text v0.3.7 // indirect
)
