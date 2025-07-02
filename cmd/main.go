package main

import (
	"be/config"
	"be/database"
	"be/internal/controllers"
	"be/internal/repositories"
	"be/internal/routes"
	"be/internal/services"
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Tải cấu hình
	cfg := config.LoadConfig()

	// Kết nối MongoDB
	mongoDB, err := database.ConnectMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close(context.Background())

	// Khởi tạo repository
	userRepo := repositories.NewUserRepository(mongoDB.UserCollection, mongoDB.Database)
	coachRepo := repositories.NewCoachRepository(mongoDB.CoachCollection, mongoDB.Database)
	athleteRepo := repositories.NewAthleteRepository(mongoDB.AthleteCollection, mongoDB.Database)
	//sportUserRepo := repositories.NewSportUserRepository(mongoDB.SportUserCollection)
	sportRepo := repositories.NewSportRepository(mongoDB.SportCollection, mongoDB.Database)
	exerciseRepo := repositories.NewExerciseRepository(mongoDB.ExerciseCollection, mongoDB.Database)
	trainingExerciseRepo := repositories.NewTrainingExerciseRepository(mongoDB.TrainingExerciseCollection)
	TrainingScheduleUserRepo := repositories.NewTrainingScheduleUserRepository(mongoDB.TrainingScheduleUserCollection)
	dailyScheduleRepo := repositories.NewDailyScheduleRepository(mongoDB.DailyScheduleCollection, mongoDB.Database)
	trainingScheduleRepo := repositories.NewTrainingScheduleRepository(mongoDB.TrainingScheduleCollection, dailyScheduleRepo, mongoDB.Database)
	notificationRepo := repositories.NewNotificationRepository(mongoDB.NotificationCollection)
	reminderRepo := repositories.NewReminderRepository(mongoDB.ReminderCollection)
	achivementRepo := repositories.NewAchievementRepository(mongoDB.AchivementCollection)
	userMatchRepo := repositories.NewUserMatchRepository(mongoDB.UserMatchCollection)
	coachCertificationRepo := repositories.NewCoachCertificationRepository(mongoDB.CoachCertificationCollection)
	feedbackRepo := repositories.NewFeedbackRepository(mongoDB.FeedbackCollection)
	groupRepo := repositories.NewGroupRepository(mongoDB.GroupCollection, mongoDB.Database)
	groupMemberRepo := repositories.NewGroupMemberRepository(mongoDB.GroupCollection)
	healthRepo := repositories.NewHealthRepository(mongoDB.HealthCollection, mongoDB.Database)
	injuryRepo := repositories.NewInjuryRepository(mongoDB.InjuryCollection)
	matchScheduleRepo := repositories.NewMatchScheduleRepository(mongoDB.MatchScheduleCollection, mongoDB.Database)
	medicalHistoryRepo := repositories.NewMedicalHistoryRepository(mongoDB.MedicalHistoryCollection)
	messageRepo := repositories.NewMessageRepository(mongoDB.MessageCollection)
	nutritionPlanRepo := repositories.NewNutritionPlanRepository(mongoDB.NutritionPlanCollection, mongoDB.Database)
	foodRepo := repositories.NewFoodRepository(mongoDB.FoodCollection, mongoDB.Database)
	//performanceRepo := repositories.NewPerformanceRepository(mongoDB.PerformanceCollection)
	//progressRepo := repositories.NewProgressRepository(mongoDB.ProgressCollection)
	teamRepo := repositories.NewTeamRepository(mongoDB.TeamCollection, mongoDB.Database)
	temaMemberRepo := repositories.NewTeamMemberRepository(mongoDB.TeamMemberCollection)
	tournamentRepo := repositories.NewTournamentRepository(mongoDB.TournamentCollection, mongoDB.Database)
	planFoodRepo := repositories.NewPlanFoodRepository(mongoDB.PlanFoodCollection)
	coachAthleteRepo := repositories.NewCoachAthleteRepository(mongoDB.CoachAthleteCollection)

	// Khởi tạo service
	userService := services.NewUserService(mongoDB.Client, userRepo, athleteRepo, coachRepo)
	coachService := services.NewCoachService(coachRepo)
	athleteService := services.NewAthleteService(athleteRepo)
	//sportUserService := services.NewSportUserService(sportUserRepo)
	sportService := services.NewSportService(sportRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	trainingExerciseService := services.NewTrainingExerciseService(trainingExerciseRepo)
	trainingScheduleService := services.NewTrainingScheduleService(trainingScheduleRepo, trainingExerciseService, trainingExerciseRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	reminderService := services.NewReminderService(reminderRepo)
	TrainingScheduleUserService := services.NewTrainingScheduleUserService(TrainingScheduleUserRepo, notificationService, reminderService, trainingScheduleRepo)
	achivementService := services.NewAchievementService(achivementRepo)
	userMatchService := services.NewUserMatchService(userMatchRepo)
	coachCertificationService := services.NewCoachCertificationService(coachCertificationRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo)
	groupService := services.NewGroupService(groupRepo)
	groupMemberService := services.NewGroupMemberService(groupMemberRepo)
	healthService := services.NewHealthService(healthRepo)
	injuryService := services.NewInjuryService(injuryRepo)
	matchScheduleService := services.NewMatchScheduleService(matchScheduleRepo)
	medicalHistoryService := services.NewMedicalHistoryService(medicalHistoryRepo)
	messageService := services.NewMessageService(messageRepo)
	planFoodService := services.NewPlanFoodService(planFoodRepo)
	nutritionPlanService := services.NewNutritionPlanService(nutritionPlanRepo, planFoodService, foodRepo)
	foodService := services.NewFoodService(foodRepo)
	//performanceService := services.NewPerformanceService(performanceRepo)
	//progressService := services.NewProgressService(progressRepo)
	teamService := services.NewTeamService(teamRepo)
	temaMemberService := services.NewTeamMemberService(temaMemberRepo)
	tournamentService := services.NewTournamentService(tournamentRepo)
	coachAthleteService := services.NewCoachAthleteService(coachAthleteRepo)
	dailyScheduleService := services.NewDailyScheduleService(dailyScheduleRepo, trainingScheduleService)
	InitCronJobs(trainingScheduleService) // Khởi tạo cron job

	// Khởi tạo controller
	userController := controllers.NewUserController(userService)
	coachController := controllers.NewCoachController(coachService)
	athleteController := controllers.NewAthleteController(athleteService)
	//sportUserController := controllers.NewSportUserController(sportUserService)
	sportController := controllers.NewSportController(sportService)
	exerciseController := controllers.NewExerciseController(exerciseService)
	trainingExerciseController := controllers.NewTrainingExerciseController(trainingExerciseService)
	TrainingScheduleUserController := controllers.NewTrainingScheduleUserController(TrainingScheduleUserService)
	trainingScheduleController := controllers.NewTrainingScheduleController(trainingScheduleService)
	imageController := controllers.NewImageController()
	notificationController := controllers.NewNotificationController(notificationService)
	reminderController := controllers.NewReminderController(reminderService)
	achivementController := controllers.NewAchievementController(achivementService)
	userMatchController := controllers.NewUserMatchController(userMatchService)
	coachCertificationController := controllers.NewCoachCertificationController(coachCertificationService)
	feedbackController := controllers.NewFeedbackController(feedbackService)
	groupController := controllers.NewGroupController(groupService)
	groupMemberController := controllers.NewGroupMemberController(groupMemberService)
	healthController := controllers.NewHealthController(healthService)
	injuryController := controllers.NewInjuryController(injuryService)
	matchScheduleController := controllers.NewMatchScheduleController(matchScheduleService)
	medicalHistoryController := controllers.NewMedicalHistoryController(medicalHistoryService)
	messageController := controllers.NewMessageController(messageService)
	nutritionPlanController := controllers.NewNutritionPlanController(nutritionPlanService)
	foodController := controllers.NewFoodController(foodService)
	//performanceController := controllers.NewPerformanceController(performanceService)
	//progressController := controllers.NewProgressController(progressService)
	teamController := controllers.NewTeamController(teamService)
	temaMemberController := controllers.NewTeamMemberController(temaMemberService)
	tournamentController := controllers.NewTournamentController(tournamentService)
	planFoodController := controllers.NewPlanFoodController(planFoodService)
	coachAthleteController := controllers.NewCoachAthleteController(coachAthleteService)
	dailyScheduleController := controllers.NewDailyScheduleController(dailyScheduleService)

	// Khởi tạo router Gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow Flutter frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Gắn routes
	routes.SetupUserRoutes(r, userController)
	routes.SetupCoachRoutes(r, coachController)
	routes.SetupAthleteRoutes(r, athleteController)
	//routes.SetupSportUserRoutes(r, sportUserController)
	routes.SetupSportRoutes(r, sportController)
	routes.SetupExerciseRoutes(r, exerciseController)
	routes.SetupTrainingExerciseRoutes(r, trainingExerciseController)
	routes.SetupTrainingScheduleUserRoutes(r, TrainingScheduleUserController)
	routes.SetupTrainingScheduleRoutes(r, trainingScheduleController)
	routes.SetupImageRoutes(r, imageController)
	routes.SetupNotificationRoutes(r, notificationController)
	routes.SetupReminderRoutes(r, reminderController)
	routes.SetupAchievementRoutes(r, achivementController)
	routes.SetupUserMatchRoutes(r, userMatchController)
	routes.SetupCoachCertificationRoutes(r, coachCertificationController)
	routes.SetupFeedbackRoutes(r, feedbackController)
	routes.SetupGroupRoutes(r, groupController)
	routes.SetupGroupMemberRoutes(r, groupMemberController)
	routes.SetupHealthRoutes(r, healthController)
	routes.SetupInjuryRoutes(r, injuryController)
	routes.SetupMatchScheduleRoutes(r, matchScheduleController)
	routes.SetupMedicalHistoryRoutes(r, medicalHistoryController)
	routes.SetupMessageRoutes(r, messageController)
	routes.SetupFoodRoutes(r, foodController)
	routes.SetupNutritionPlanRoutes(r, nutritionPlanController)
	//routes.SetupPerformanceRoutes(r, performanceController)
	//routes.SetupProgressRoutes(r, progressController)
	routes.SetupTeamMemberRoutes(r, temaMemberController)
	routes.SetupTeamRoutes(r, teamController)
	routes.SetupTournamentRoutes(r, tournamentController)
	routes.SetupPlanFoodRoutes(r, planFoodController)
	routes.SetupCoachAthleteRoutes(r, coachAthleteController)
	routes.SetupDailyScheduleRoutes(r, dailyScheduleController)

	// Chạy server
	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
