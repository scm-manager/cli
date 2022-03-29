package pkg

type Configuration struct {
	ServerUrl string
	Username  string
	ApiKey    string `json:"-"`
}
