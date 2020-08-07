module github.com/saltbo/authcar

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/saltbo/gopkg v0.0.0-20200807091256-33dd5deb87d0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/storyicon/grbac v0.0.0-20200224041032-a0461737df7e
	gopkg.in/resty.v1 v1.12.0
)

//replace github.com/saltbo/gopkg => /opt/works/gopkg
