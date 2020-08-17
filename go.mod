module github.com/saltbo/moreu

go 1.14

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/magiconair/properties v1.8.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/saltbo/gopkg v0.0.0-20200817013558-116ba552a859
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/storyicon/grbac v0.0.0-20200224041032-a0461737df7e
	github.com/swaggo/swag v1.6.7
	gopkg.in/resty.v1 v1.12.0 // indirect
)

//replace github.com/saltbo/gopkg => /opt/works/gopkg
