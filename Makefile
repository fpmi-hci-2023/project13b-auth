proto:
	protoc -Iapi/auth/v1 --go_out=api/auth/v1/gen \
 		--go_opt=module=github.com/fpmi-hci-2023/project13b-auth/api/auth/v1 \
		--go-grpc_out=api/auth/v1/gen \
		--go-grpc_opt=module=github.com/fpmi-hci-2023/project13b-auth/api/auth/v1 \
		api/auth/v1/auth.proto

migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="host=contabo.richardhere.dev port=30570 user=hci password=E8n930d1SdXeOmZMLnQtKeGVegqfNXql05xyLi1I6vSS12KLa3Tiw3DgI2FWTf3m dbname=auth sslmode=require" goose up
