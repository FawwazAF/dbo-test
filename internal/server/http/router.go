package http

func (s *Server) registerHandler() {
	// authWrapper := s.authMiddleware
	s.router.GET("/", s.handler.Index.HandlerIndex)

	s.router.GET("/v1/customer/:id", s.handler.Customer.HandlerGetCustomerByID)
	s.router.POST("/v1/customer", s.handler.Customer.HandlerAddCustomer)
	s.router.PUT("/v1/customer", s.handler.Customer.HandlerUpdateCustomer)
	s.router.DELETE("/v1/customer/:id", s.handler.Customer.HandlerDeleteCustomer)
	s.router.GET("/v1/customer", s.handler.Customer.HandlerSearchCustomer)
}
