package main

import (
	"bolado-stack/libs/database"
	"bolado-stack/src/ports"
	"bolado-stack/src/repositories"
	"bolado-stack/src/services"
)

func main() {
	mongoConnection := database.NewMongoConnection(database.MongoConfig{})

	// repositories init
	repositories := repositories.Container{
		UserMongoDBRepository: repositories.NewUserMongoDBRepository(mongoConnection),
	}

	// services init
	services := services.Container{
		UserService: services.NewUserService(repositories),
	}

	httpPortConfig := ports.HTTPPortConfig{
		Services:   services,
		Debug:      true,
		HideBanner: false,
	}

	ports.SetupHTTPServer(httpPortConfig)
}
