package entpb

//go:generate protoc -I=.. --go_out=../gen --go-grpc_out=../gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --entgrpc_out=../gen --entgrpc_opt=paths=source_relative,schema_path=../../ent/schema entpb/entpb.proto entpb/ext.proto
