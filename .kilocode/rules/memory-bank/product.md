# Product Description

## What This Project Is
An e-wallet application that provides digital banking services with a focus on savings accounts and interest calculation. Users can create accounts, manage their profiles, and earn interest on their savings through automated daily calculations.

## Problems It Solves
- **Digital Banking Access**: Provides users with a digital alternative to traditional banking for savings and payments
- **Interest Earning**: Automates interest calculation and application for savings accounts with tiered rates based on balance and account age
- **Account Management**: Allows users to manage multiple account types (payment and savings) with proper limits and validation
- **Transaction Tracking**: Maintains complete history of all transactions and interest payments
- **Profile Management**: Enables users to maintain personal information and banking profiles

## How It Should Work
1. **User Registration**: Users register with email/password and create a profile
2. **Account Creation**: Users can create one payment account and multiple savings accounts (with limits)
3. **Daily Interest Calculation**: A background worker calculates and applies interest to flexible savings accounts based on tiered rates
4. **Transaction Recording**: All balance changes are recorded as transactions
5. **Interest History**: Interest payments are tracked separately for transparency

## User Experience Goals
- **Secure**: JWT authentication ensures secure access
- **Reliable**: Hexagonal architecture provides maintainable and testable code
- **Scalable**: PostgreSQL and proper indexing support growth
- **Transparent**: Clear transaction and interest history for users
- **User-Friendly**: RESTful API with proper error handling and validation