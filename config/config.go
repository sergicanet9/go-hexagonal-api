package config

const (
	// DbConnectionString provides the mongo connection string
	DbConnectionString = "mongodb+srv://admin:admin@cluster0.qy8ev.mongodb.net/go-mongo-restapi?retryWrites=true&w=majority"
	//DbName is the name of the used database
	DbName = "go-mongo-restapi"
	// APIPort is the port used when running the app
	APIPort = 8080
	//JWTSecret is the secret key for creating JWT Tokens
	JWTSecret = "CTeemck6Gg"
)
