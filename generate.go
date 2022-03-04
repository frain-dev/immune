package immune

//go:generate mockgen --source callback.go --destination mocks/callback.go -package mocks
//go:generate mockgen --source database/truncator.go --destination mocks/truncator.go -package mocks
