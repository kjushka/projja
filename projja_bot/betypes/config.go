package betypes

import (
	"fmt"
)

var (
	BotToken           = "1436657093:AAEp-Vsd91oOWWjfDOhcfn9bNc5wmVzj0yw"
	BotExternalAddress = "34.92.36.31"
	BotExternalPort    = "443"
	BotInternalAddress = "0.0.0.0"
	BotInternalPort    = "5000"
	TelegramUrl        = "https://api.telegram.org/bot"
	MySqlAddress 			 = "localhost"
	MySqlPort 				 = "8080"	
)

// Агрумент http - вернет путь по http протоколу
// Агрумент https - вернет путь по https протоколу
func GetPathToMySQl(protType string) string {
	if(protType == "http") {
		return fmt.Sprintf("http://%s:%s/", MySqlAddress, MySqlPort);
	}
	
	return fmt.Sprintf("https://%s:%s/", MySqlAddress, MySqlPort);
}


