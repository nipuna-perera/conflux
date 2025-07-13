# Development User & Authentication

This project includes development utilities to streamline the development workflow by bypassing login/signup when needed.

## 🔑 Default Development User

**Email**: `dev@conflux.local`  
**Password**: `password123`

This user is automatically created when you run `make dev` in development environment.

## 📋 Development Options

### Option 1: Use Regular Login/Signup
Access the app normally at `http://localhost:3000` and create accounts as needed.

### Option 2: Get Development JWT Token
Get a pre-authenticated JWT token without going through login flow:

```bash
# Get development token
curl -X POST http://localhost:8080/dev/token

# Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "email": "dev@conflux.local",
    "first_name": "Dev", 
    "last_name": "User"
  },
  "instructions": "Use this token in Authorization header: Bearer eyJhbGciOiJIUzI1NiIs..."
}
```

### Option 3: Ensure Dev User Exists
If the migration didn't create the dev user, manually create it:

```bash
curl -X POST http://localhost:8080/dev/user
```

## 🛠️ Using the Development Token

### In API Requests
```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
     http://localhost:8080/api/protected-endpoint
```

### In Frontend Development
Store the token in localStorage and use it for authenticated requests:

```javascript
// Get the token
const response = await fetch('http://localhost:8080/dev/token', {
  method: 'POST'
});
const data = await response.json();

// Store for use
localStorage.setItem('auth_token', data.token);

// Use in requests
fetch('/api/protected', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
  }
});
```

## 🔒 Security Considerations

### ✅ Good Practices Implemented:
- **Environment-gated**: Dev endpoints only work when `ENVIRONMENT=development`
- **Predictable credentials**: Makes development consistent across team
- **Separate from production**: No impact on production security
- **Well-documented**: Clear usage instructions
- **Fallback available**: Regular auth flow still works

### 🚫 What We Avoid:
- **Never in production**: Dev endpoints return 403 in non-dev environments
- **No hardcoded secrets**: Uses same JWT secret as regular auth
- **No backdoors**: Dev user follows same validation as regular users

## 🏗️ Industry Examples

This pattern is used by major companies:

- **GitHub**: Uses Octocat user for demos and testing
- **Stripe**: Provides test API keys and sample data  
- **Auth0**: Has development tenants with sample users
- **Firebase**: Includes demo projects with pre-configured auth

## 🔄 Development Workflow

```bash
# 1. Start development environment
make dev

# 2. Get dev token (optional)
curl -X POST http://localhost:8080/dev/token

# 3. Start developing with authentication already handled
# or use regular signup/login flow
```

## 🎯 Benefits

- **Faster iteration**: No repeated login steps
- **Testing authenticated flows**: Easy to test protected features
- **Team consistency**: Everyone uses same baseline data
- **CI/CD friendly**: Automated tests can get tokens instantly
- **Demo ready**: Always have a working authenticated state

This approach follows industry best practices for development experience while maintaining production security.
