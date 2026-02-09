variable "GO_VERSION" {
  default = "1.25.6"
}

variable "GAME_IMAGE" {
  default = "docker.io/louisbranch/fracturing.space-game:dev"
}

variable "MCP_IMAGE" {
  default = "docker.io/louisbranch/fracturing.space-mcp:dev"
}

variable "ADMIN_IMAGE" {
  default = "docker.io/louisbranch/fracturing.space-admin:dev"
}

variable "AUTH_IMAGE" {
  default = "docker.io/louisbranch/fracturing.space-auth:dev"
}

group "default" {
  targets = ["game", "mcp", "admin", "auth"]
}

target "base" {
  context    = "."
  dockerfile = "Dockerfile"
  args = {
    GO_VERSION = "${GO_VERSION}"
  }
}

target "game" {
  inherits = ["base"]
  target   = "game"
  tags     = ["${GAME_IMAGE}"]
}

target "mcp" {
  inherits = ["base"]
  target   = "mcp"
  tags     = ["${MCP_IMAGE}"]
}

target "admin" {
  inherits = ["base"]
  target   = "admin"
  tags     = ["${ADMIN_IMAGE}"]
}

target "auth" {
  inherits = ["base"]
  target   = "auth"
  tags     = ["${AUTH_IMAGE}"]
}
