# protoc calculator.proto --go_out=../server --go-grpc_out=../server
# protoc calculator.proto --go_out=../client --go-grpc_out=../client

# ถ้ามีไฟล์ proto หลายอันให้ใช้คำสั่ง 
protoc *.proto --go_out=../server --go-grpc_out=../server
protoc *.proto --go_out=../client --go-grpc_out=../client
# โปรแกรมจะ gen file proto ให้ทั้งหมด