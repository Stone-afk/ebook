package service

import (
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/service/sms"
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
	//tplId string
}

type NamedArg struct {
	Val  string
	Name string
}
