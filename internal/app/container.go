package app

import (
	"go-split/internal/domain/repository"
	"go-split/internal/domain/usecase"
	"go-split/internal/infrastructure/database"
	infrastructureRepository "go-split/internal/infrastructure/repository"
	"go-split/internal/interface/http"
	"go-split/internal/interface/http/handler"
	"go-split/pkg/config"
	"go-split/pkg/libs/helper"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	Config             *config.Config
	Router             *gin.Engine
	MongoDB            *mongo.Database
	CloudinaryUploader *helper.CloudinaryUploader
	Repository         *Repository
	UseCase            *UseCase
	Handler            *Handler
}

type Repository struct {
	UserRepository         repository.UserRepository
	GroupRepository        repository.GroupRepository
	ExpenseRepository      repository.ExpenseRepository
	ExpenseSplitRepository repository.ExpenseSplitRepository
}

type UseCase struct {
	UserUseCase    usecase.UserUseCase
	GroupUseCase   usecase.GroupUseCase
	ExpenseUseCase usecase.ExpenseUseCase
}

type Handler struct {
	UserHandler    *handler.UserHandler
	GroupHandler   *handler.GroupHandler
	ExpenseHandler *handler.ExpenseHandler
}

func NewContainer() (*Container, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	cloudinaryUploader, err := helper.NewCloudinaryUploader(cfg.Cloudinary.URL)
	if err != nil {
		return nil, err
	}

	c := &Container{
		Config:             cfg,
		CloudinaryUploader: cloudinaryUploader,
		Repository:         &Repository{},
		UseCase:            &UseCase{},
		Handler:            &Handler{},
	}

	if err := c.initDatabase(); err != nil {
		return nil, err
	}

	c.initRouter()

	c.initRepositories()

	c.initUseCases()

	c.initHandlers()

	c.setupRouter()

	return c, nil
}

func (c *Container) initDatabase() error {
	mongoDB, err := database.NewMongoConnection(c.Config.MongoDB)
	if err != nil {
		return err
	}

	c.MongoDB = mongoDB
	log.Println("MongoDB connection established successfully")

	return nil
}

func (c *Container) initRouter() {
	c.Router = gin.Default()

	// Configure CORS
	c.Router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://react-split.vercel.app",
		}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Refresh-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func (c *Container) setupRouter() {
	http.SetupRouter(c.Router, c.Handler.UserHandler, c.Handler.GroupHandler, c.Handler.ExpenseHandler)
}

func (c *Container) initRepositories() {
	c.Repository.UserRepository = infrastructureRepository.NewUserRepositoryMongo(c.MongoDB.Collection("users"))
	c.Repository.GroupRepository = infrastructureRepository.NewGroupRepositoryMongo(c.MongoDB.Collection("groups"))
	c.Repository.ExpenseRepository = infrastructureRepository.NewExpenseRepositoryMongo(c.MongoDB.Collection("expenses"))
	c.Repository.ExpenseSplitRepository = infrastructureRepository.NewExpenseSplitRepository(c.MongoDB.Collection("expense_splits"))
}

func (c *Container) initUseCases() {
	c.UseCase.UserUseCase = usecase.NewUserUseCase(c.Repository.UserRepository, c.CloudinaryUploader)
	c.UseCase.GroupUseCase = usecase.NewGroupUseCase(c.Repository.GroupRepository, c.Repository.UserRepository, c.Repository.ExpenseRepository, c.Repository.ExpenseSplitRepository, c.CloudinaryUploader)
	c.UseCase.ExpenseUseCase = usecase.NewExpenseUseCase(c.Repository.ExpenseRepository, c.Repository.ExpenseSplitRepository, c.Repository.GroupRepository, c.Repository.UserRepository, c.CloudinaryUploader)
}

func (c *Container) initHandlers() {
	c.Handler.UserHandler = handler.NewUserHandler(c.UseCase.UserUseCase)
	c.Handler.GroupHandler = handler.NewGroupHandler(c.UseCase.GroupUseCase)
	c.Handler.ExpenseHandler = handler.NewExpenseHandler(c.UseCase.ExpenseUseCase)
}
