package config


const configFileName = ".gatorconfig.json"

// todo: export aconfig struct  representing json structure with tags

// export read function  reads the json file at ~/.gatorconfig.json returns Config struct
// --should read from the home directory, then decode the json string into a new config struct
// --use "os.UserHomeDir" to get location of home

// export a "SetUser" method on the "Config" struct that writes the config struct to the  JSON file
// after setting "current_user_name" field


//helper functions:
func getConfigFilePath() (string, error){

}

func write(cfg Config) error {

}
