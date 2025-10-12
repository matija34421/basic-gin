package main

// @title       Basic Gin Service API
// @version     1.0
// @description API za klijente, račune i transakcije.
// @BasePath    /api/v1

// @tag.name        clients
// @tag.description Operacije nad klijentima
// @tag.name        accounts
// @tag.description Operacije nad računima
// @tag.name        transactions
// @tag.description Operacije nad transakcijama

// HMAC preko jednog headera (najjednostavnija varijanta)
// @securityDefinitions.apikey HmacSignature
// @in header
// @name X-Signature
