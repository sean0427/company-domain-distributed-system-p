mock_gen: 
	mockgen -package=mock --source=postgressql/rdb.go > mock/postgressql_client_mock.go
	mockgen -package=mock --source=company.go > mock/repo_mock.go
	mockgen -package=mock --source=grpc/grpc.go > mock/service_mock.go

protoco_gen: 
	protoc proto/*.proto --go_out=${PWD} --go-grpc_out=${PWD} --experimental_allow_proto3_optional

build_push_to_kind:
	docker build . -t compnay-domain
	kind load docker-image compnay-domain --name micro-service
