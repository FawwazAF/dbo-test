package http

func (s *Server) registerHandler() {
	authWrapper := s.authMiddleware
	s.router.GET("/", s.handler.Index.HandlerIndex)
	s.router.POST("/login", s.handler.Login.HandlerLogin)
	s.router.GET("/login", authWrapper(), s.handler.Login.HandlerLoginInfo)

	s.router.GET("/v1/customer/:id", s.handler.Customer.HandlerGetCustomerByID)
	s.router.POST("/v1/customer", s.handler.Customer.HandlerAddCustomer)
	s.router.PUT("/v1/customer", s.handler.Customer.HandlerUpdateCustomer)
	s.router.DELETE("/v1/customer/:id", authWrapper(SuperAdminAccess), s.handler.Customer.HandlerDeleteCustomer)
	s.router.GET("/v1/customer", s.handler.Customer.HandlerSearchCustomer)

	s.router.GET("/v1/order/:order_id", authWrapper(), s.handler.Order.HandlerGetOrderDetail)
	s.router.POST("/v1/order", authWrapper(), s.handler.Order.HandlerCreateOrder)
	s.router.DELETE("/v1/order/:order_id", authWrapper(), s.handler.Order.HandlerDeleteOrder)
	s.router.PUT("/v1/order", authWrapper(), s.handler.Order.HandlerUpdateOrder)
	s.router.GET("/v1/order", authWrapper(), s.handler.Order.HandlerSearchOrder)

	s.router.GET("/v1/product", s.handler.Product.HandlerGetAllProduct)
}
