module github.com/NomanSalhab/golang_b_n_b_training_project

go 1.19

replace github.com/nomansalhab/golang_b_n_b_training_project/cmd/web/ => ./cmd/web/

require (
	github.com/alexedwards/scs/v2 v2.5.0
	github.com/go-chi/chi v1.5.4
	github.com/justinas/nosurf v1.1.1
)

require github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
