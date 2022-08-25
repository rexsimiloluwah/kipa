generate-repository-mocks: 
	mockgen --source repository/main.go --destination mocks/repository.go -package mocks