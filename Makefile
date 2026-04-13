gen-user:
	@protoc \
		--proto_path=proto "proto/user.proto" \
		--go_out=common/gen-proto/users --go_opt=paths=source_relative \
  	--go-grpc_out=common/gen-proto/users --go-grpc_opt=paths=source_relative

gen-auth:
	@protoc \
		--proto_path=proto "proto/auth.proto" \
		--go_out=common/gen-proto/auth --go_opt=paths=source_relative \
  	--go-grpc_out=common/gen-proto/auth --go-grpc_opt=paths=source_relative

gen-product:
	@protoc \
		--proto_path=proto "proto/product.proto" \
		--go_out=common/gen-proto/products --go_opt=paths=source_relative \
  	--go-grpc_out=common/gen-proto/products --go-grpc_opt=paths=source_relative

gen-order:
	@protoc \
		--proto_path=proto "proto/order.proto" \
		--go_out=common/gen-proto/orders --go_opt=paths=source_relative \
  	--go-grpc_out=common/gen-proto/orders --go-grpc_opt=paths=source_relative

gen-otp:
	@protoc \
		--proto_path=proto "proto/otp.proto" \
		--go_out=common/gen-proto/otp --go_opt=paths=source_relative \
  	--go-grpc_out=common/gen-proto/otp --go-grpc_opt=paths=source_relative