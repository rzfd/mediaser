package configs

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Auth     AuthConfig
	Payment  PaymentConfig
	Midtrans MidtransConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

type AuthConfig struct {
	JWTSecret    string
	TokenExpiry  int
}

type PaymentConfig struct {
	PaypalClientID     string
	PaypalClientSecret string
	StripeSecretKey    string
	WebhookSecret      string
}

type MidtransConfig struct {
	MerchantID    string
	ClientKey     string
	ServerKey     string
	Environment   string // sandbox or production
	WebhookSecret string
}

// LoadConfig loads configuration from config file and environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	// Override with environment variables if they exist
	if os.Getenv("DB_HOST") != "" {
		config.DB.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		config.DB.Port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_USERNAME") != "" {
		config.DB.Username = os.Getenv("DB_USERNAME")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		config.DB.Password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_NAME") != "" {
		config.DB.Name = os.Getenv("DB_NAME")
	}
	if os.Getenv("SERVER_PORT") != "" {
		config.Server.Port = os.Getenv("SERVER_PORT")
	}
	if os.Getenv("JWT_SECRET") != "" {
		config.Auth.JWTSecret = os.Getenv("JWT_SECRET")
	}
	
	// Midtrans environment variables
	if os.Getenv("MIDTRANS_MERCHANT_ID") != "" {
		config.Midtrans.MerchantID = os.Getenv("MIDTRANS_MERCHANT_ID")
	}
	if os.Getenv("MIDTRANS_CLIENT_KEY") != "" {
		config.Midtrans.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	}
	if os.Getenv("MIDTRANS_SERVER_KEY") != "" {
		config.Midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	}
	if os.Getenv("MIDTRANS_ENVIRONMENT") != "" {
		config.Midtrans.Environment = os.Getenv("MIDTRANS_ENVIRONMENT")
	}

	return config, nil
} 