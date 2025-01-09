package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abedsully/golang-microservice/api/server"
	"github.com/abedsully/golang-microservice/api/storer"
	"github.com/abedsully/golang-microservice/token"
	"github.com/abedsully/golang-microservice/util"
	"github.com/go-chi/chi"
)

type handler struct {
	ctx    context.Context
	server *server.Server
	tokenMaker *token.JWTMaker
}

func NewHandler(server *server.Server, secretKey string) *handler {
	return &handler{
		ctx:    context.Background(),
		server: server,
		tokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	product, err := h.server.CreateProduct(h.ctx, toStorerProduct(p))

	if err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)

	w.Header().Set("Contet-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}

	product, err := h.server.GetProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.server.GetAllProducts(h.ctx)

	if err != nil {
		http.Error(w, "error getting all products", http.StatusInternalServerError)
		return
	}

	var res []ProductRes
	for _, p := range products {
		res = append(res, toProductRes(&p))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}

	var p ProductReq

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	product, err := h.server.GetProduct(h.ctx, i)

	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	patchProductReq(product, p)

	updated, err := h.server.UpdateProduct(h.ctx, product)
	if err != nil {
		http.Error(w, "error updating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(updated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}

	if err := h.server.DeleteProduct(h.ctx, i); err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func toStorerProduct(p ProductReq) *storer.Product {
	return &storer.Product{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func toProductRes(p *storer.Product) ProductRes {
	return ProductRes{
		ID:           p.ID,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func patchProductReq(product *storer.Product, p ProductReq) {
	if p.Name != "" {
		product.Name = p.Name
	}

	if p.Image != "" {
		product.Image = p.Image
	}

	if p.Category != "" {
		product.Category = p.Category
	}

	if p.Description != "" {
		product.Description = p.Description
	}

	if p.Rating != 0 {
		product.Rating = p.Rating
	}

	if p.NumReviews != 0 {
		product.NumReviews = p.NumReviews
	}

	if p.Price != 0 {
		product.Price = p.Price
	}

	if p.CountInStock != 0 {
		product.CountInStock = p.CountInStock
	}

	product.UpdatedAt = toTimePtr(time.Now())
}

func toTimePtr(t time.Time) *time.Time {
	return &t
}

func (h *handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var o OrderReq

	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "error decoding json body", http.StatusBadRequest)
		return
	}

	created, err := h.server.CreateOrder(h.ctx, toStorerOrder(o))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	res := toOrderRes(created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

	
}

func (h *handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	i, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
	}

	order, err := h.server.GetOrder(h.ctx, i)

	if err != nil {
		http.Error(w, "error getting order", http.StatusInternalServerError)
	}

	res := toOrderRes(order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func(h *handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.server.GetAllOrders(h.ctx)

	if err != nil {
		http.Error(w, "error getting all orders", http.StatusInternalServerError)
		return
	}

	var res []OrderRes
	for _, o := range orders {
		res = append(res, toOrderRes(&o))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func(h *handler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}

	if err := h.server.DeleteOrder(h.ctx, i); err != nil {
		http.Error(w, "error deleting order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toStorerOrder(o OrderReq) *storer.Order {
	return &storer.Order{
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Items:         toStorerOrderItems(o.Items),
	}
}

func toStorerOrderItems(items []OrderItem) []storer.OrderItem {
	var res []storer.OrderItem
	for _, i := range items {
		res = append(res, storer.OrderItem{
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductID: i.ProductID,
		})
	}
	return res
}

func toOrderRes(o *storer.Order) OrderRes {
	return OrderRes{
		ID:            o.ID,
		Items:         toOrderItems(o.Items),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
}

func toOrderItems(items []storer.OrderItem) []OrderItem {
	var res []OrderItem
	for _, i := range items {
		res = append(res, OrderItem{
			Name:      i.Name,
			Quantity:  i.Quantity,
			Image:     i.Image,
			Price:     i.Price,
			ProductID: i.ProductID,
		})
	}
	return res
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding json body", http.StatusBadRequest)
		return
	}

	hashed, err := util.HashPassword(u.Password)

	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	u.Password = hashed

	created, err := h.server.CreateUser(h.ctx, toStorerUser(&u))

	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	res := toUserRes(created)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func toStorerUser(u *UserReq) *storer.User {
	return &storer.User {
		Name: u.Name,
		Email: u.Email,
		Password: u.Password,
		IsAdmin: u.IsAdmin,
	}
}

func toUserRes(u *storer.User) UserRes {
	return UserRes{
		Name: u.Name,
		Email: u.Email,
		IsAdmin: u.IsAdmin,
	}
}

func patchUserReq(u *storer.User, p UserReq) {
	if p.Name != "" {
		u.Name = p.Name
	}

	if p.Email != "" {
		u.Email = p.Email
	}

	if p.Password != "" {
		hashedPassword, err := util.HashPassword(p.Password)
		if err != nil {
			panic(err)
		}

		u.Password = hashedPassword
	}

	if p.IsAdmin {
		u.IsAdmin = p.IsAdmin
	}

	u.UpdatedAt = toTimePtr(time.Now())
} 

func (h *handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.server.GetAllUsers(h.ctx)

	if err != nil {
		http.Error(w, "error getting all users", http.StatusInternalServerError)
		return
	}

	var res AllUsers

	for _, u := range users {
		res.Users = append(res.Users, toUserRes(&u))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *handler) updateUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding json body", http.StatusBadRequest)
		return
	}

	user, err := h.server.GetUser(h.ctx, u.Email)

	fmt.Println(user.ID)

	if err != nil {
		http.Error(w, "error getting user", http.StatusBadRequest)
		return
	}

	patchUserReq(user, u)

	updated, err := h.server.UpdateUser(h.ctx, user)

	if err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	res := toUserRes(updated)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	i, err := strconv.ParseInt(id, 10, 64)
	
	if err != nil {
		http.Error(w, "error parsing int", http.StatusBadRequest)
		return
	}

	if err := h.server.DeleteUser(h.ctx, i); err != nil {
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var u LoginUserReq

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding json body", http.StatusBadRequest)
		return
	}

	user, err := h.server.GetUser(h.ctx, u.Email)

	if err != nil {
		http.Error(w, "error getting user", http.StatusBadRequest)
		return
	}

	err = util.CheckPassword(u.Password, user.Password)

	if err != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := h.tokenMaker.CreateToken(user.ID, user.Email, user.IsAdmin, 15 * time.Minute)

	if err != nil {
		http.Error(w, "error creating access token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaim, err := h.tokenMaker.CreateToken(user.ID, user.Email, user.IsAdmin, 24 * time.Hour)

	if err != nil {
		http.Error(w, "error creating refresh token", http.StatusInternalServerError)
		return
	}

	session, err := h.server.CreateSession(h.ctx, &storer.Session{
		ID: refreshClaim.RegisteredClaims.ID,
		UserEmail: user.Email,
		RefreshToken: refreshToken,
		IsRevoked: false,
		ExpiresAt: refreshClaim.RegisteredClaims.ExpiresAt.Time,
	})

	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	res := LoginUserRes{
		SessionID: session.ID,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaim.RegisteredClaims.ExpiresAt.Time,
		User: toUserRes(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) logoutUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing session id", http.StatusBadRequest)
		return
	}

	err := h.server.DeleteSession(h.ctx, id)

	if err != nil {
		http.Error(w, "error deleting session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req RenewAccessTokenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding json body", http.StatusBadRequest)
		return
	}

	refreshClaim, err := h.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		http.Error(w, "error verifying refresh token", http.StatusUnauthorized)
		return
	}

	session, err := h.server.GetSession(h.ctx, refreshClaim.RegisteredClaims.ID)

	if err != nil {
		http.Error(w, "error getting session", http.StatusInternalServerError)
		return
	}

	if session.IsRevoked{
		http.Error(w, "session has been revoked", http.StatusUnauthorized)
		return
	}

	if session.UserEmail != refreshClaim.Email {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := h.tokenMaker.CreateToken(refreshClaim.ID, refreshClaim.Email, refreshClaim.IsAdmin, 15 * time.Minute)

	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	res := RenewAccessTokenRes {
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func(h *handler) revokeSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing session ID", http.StatusBadRequest)
		return
	}

	err := h.server.RevokeSession(h.ctx, id)

	if err != nil {
		http.Error(w, "error revoking session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

