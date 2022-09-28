BUILD=build
CPPFLAGS = -I/usr/local/include/opencv4
LDFLAGS = -lopencv_core -lopencv_calib3d -lopencv_imgproc

all: clean build

build: go.sum
	@echo "Building ..."
	@CGO_CPPFLAGS="$(CPPFLAGS)" CGO_LDFLAGS="$(LDFLAGS)" go build -mod=readonly -o $(BUILD)/arcface-api

go.sum: go.mod
	@echo "Ensure dependencies have not been modified"
	@GO111MODULE=on go mod verify

clean:
	@echo "Clean old built"
	@rm -rf $(BUILD)
	@go clean
