# Go
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=get5
MIGRATION_BINARY_NAME=migrate
ORIGINAL_MIGRATION_BINARY_NAME=migrate_prepare
DIST_DIR=build
SERVER_DIR=server
WEB_DIR=web
#GAME_DIR=game_plugin
OS_Linux=linux
OS_Windows=windows
OS_Mac=darwin
ARCH_386=386
ARCH_AMD64=amd64

# Replacing "RM" command for Windows PowerShell.
RM = rm -rf
ifeq ($(OS),Windows_NT)
    RM = Remove-Item -Recurse -Force
endif

# Replacing "MKDIR" command for Windows PowerShell.
MKDIR = mkdir -p
ifeq ($(OS),Windows_NT)
    MKDIR = New-Item -ItemType Directory
endif

# Replacing "CP" command for Windows PowerShell.
CP = cp -R
ifeq ($(OS),Windows_NT)
	CP = powershell -Command Copy-Item -Recurse -Force
endif

# Replacing "GOPATH" command for Windows PowerShell.
GOPATHDIR = $GOPATH
ifeq ($(OS),Windows_NT)
    GOPATHDIR = $$env:GOPATH
endif

.DEFAULT_GOAL := build-all

test:
	$(GOTEST) -v ./...
clean:
	@$(GOCLEAN)
	-@$(RM) $(DIST_DIR)/*
deps: deps-web deps-go
	@git submodule update
deps-web:
	@yarn global add @vue/cli
	@cd ./web && yarn
build-all: clean build-prepare build-web
	@cd ./server && gox \
	-os="$(OS_Windows) $(OS_Mac) $(OS_Linux)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../$(DIST_DIR)/$(BINARY_NAME)_{{.OS}}_{{.Arch}}/$(BINARY_NAME)"
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_386)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_386)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_386)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_386)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_386)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_386)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/static

	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_386)/config.ini.template
	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/config.ini.template
	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_386)/config.ini.template
	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)/config.ini.template
	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_386)/config.ini.template
	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/config.ini.template
	
build-prepare:
	@cd ./server && $(GOGET) github.com/mitchellh/gox \
	github.com/konsorten/go-windows-terminal-sequences \
	github.com/FlowingSPDG/get5-web-go/server
	-@$(RM) ./$(DIST_DIR)/*/static
build-linux: build-prepare build-web build-linux-server-only
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_386)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_386)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/static
build-linux-server-only: build-prepare 
	@cd ./server && gox \
	-os="$(OS_Linux)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../$(DIST_DIR)/$(BINARY_NAME)_{{.OS}}_{{.Arch}}/$(BINARY_NAME)"
	@cp ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_386)/config.ini.template
	@cp ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/config.ini.template
build-windows: build-prepare build-web build-windows-server-only
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_386)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_386)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)/static
build-windows-server-only: build-prepare
	@cd ./server && gox \
	-os="$(OS_Windows)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../$(DIST_DIR)/$(BINARY_NAME)_{{.OS}}_{{.Arch}}/$(BINARY_NAME)"
	$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_386)
	$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)
build-mac: build-prepare build-web build-mac-server-only
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_386)/static
	@$(MKDIR) ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_386)/static
	@$(CP) ./web/dist/* ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/static
build-mac-server-only: build-prepare
	@cd ./server && gox \
	-os="$(OS_Mac)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../$(DIST_DIR)/$(BINARY_NAME)_{{.OS}}_{{.Arch}}/$(BINARY_NAME)"
	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_386)/config.ini.template
	@$(CP) ./server/config.ini.template ./$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/config.ini.template
build-web:
	@cd ./web && yarn run build
build-migration-all: build-migration-linux build-migration-windows build-migration-mac
build-migration-linux:
	@cd ./migration/migrate && gox \
	-os="$(OS_Linux)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../../$(DIST_DIR)/$(MIGRATION_BINARY_NAME)_{{.OS}}_{{.Arch}}/$(MIGRATION_BINARY_NAME)"
	@cd ./migration/migrate_from_original && gox \
	-os="$(OS_Linux)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../../$(DIST_DIR)/$(ORIGINAL_MIGRATION_BINARY_NAME)_{{.OS}}_{{.Arch}}/$(ORIGINAL_MIGRATION_BINARY_NAME)"
build-migration-windows:
	@cd ./migration/migrate && gox \
	-os="$(OS_Windows)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../../$(DIST_DIR)/$(MIGRATION_BINARY_NAME)_{{.OS}}_{{.Arch}}/$(MIGRATION_BINARY_NAME)"
	@cd ./migration/migrate_from_original && gox \
	-os="$(OS_Windows)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../../$(DIST_DIR)/$(ORIGINAL_MIGRATION_BINARY_NAME)_{{.OS}}_{{.Arch}}/$(ORIGINAL_MIGRATION_BINARY_NAME)"
build-migration-mac:
	@cd ./migration/migrate && gox \
	-os="$(OS_Mac)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../../$(DIST_DIR)/$(MIGRATION_BINARY_NAME)_{{.OS}}_{{.Arch}}/$(MIGRATION_BINARY_NAME)"
	@cd ./migration/migrate_from_original && gox \
	-os="$(OS_Mac)" \
	-arch="$(ARCH_386) $(ARCH_AMD64)" \
	--output "../../$(DIST_DIR)/$(ORIGINAL_MIGRATION_BINARY_NAME)_{{.OS}}_{{.Arch}}/$(ORIGINAL_MIGRATION_BINARY_NAME)"
build-protobuf:
	@cd ./server/proto && protoc --go_out=plugins=grpc:../src/grpc/proto ./get5-web-go.proto