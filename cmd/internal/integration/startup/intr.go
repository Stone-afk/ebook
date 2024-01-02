package startup

import (
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"ebook/cmd/interactive/service"
	"ebook/cmd/internal/handler/client/interactive"
)

func InitInteractiveClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	return interactive.NewServiceAdapter(svc)
}
