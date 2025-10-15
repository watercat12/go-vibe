# E-Wallet Project Brief

## Project Overview
This is an e-wallet application built in Go using Hexagonal Architecture. It provides digital wallet functionality with user management, account creation, and savings features including interest calculation.

## Core Purpose
The application serves as a digital banking platform where users can:
- Create and manage user accounts with authentication
- Create payment accounts and savings accounts (flexible and fixed-term)
- Calculate and apply daily interest on flexible savings accounts
- Track transaction history and interest payments
- Manage user profiles with personal information

## Key Features
- User registration and JWT-based authentication
- Multiple account types: payment accounts and savings accounts
- Interest calculation for savings accounts with tiered rates
- Transaction recording and balance tracking
- Profile management
- RESTful API with Echo framework

## Technical Foundation
- Built with Go 1.25
- PostgreSQL database with GORM ORM
- Echo web framework
- Hexagonal architecture for clean separation of concerns
- Comprehensive testing with Go's testing framework and testify
- Database migrations with golang-migrate