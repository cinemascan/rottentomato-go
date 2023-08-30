package server

func Init(port string) {
	router := NewRouter()
	router.Run(":" + port)
}
