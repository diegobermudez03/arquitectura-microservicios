# Test script for Cards Microservice API

echo "Testing Cards Microservice API"
echo "=============================="

# Test 1: Register a user
echo "1. Registering a user..."
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"John\",
    \"lastname\": \"Doe\",
    \"birth_date\": \"2002-03-25\",
    \"country_code\": \"US\"
  }"

echo -e "\n\n"

# Test 2: Health check
echo "2. Health check..."
curl -X GET http://localhost:8080/health

echo -e "\n\n"

# Test 3: Issue a card (replace USER_TOKEN with actual token from step 1)
echo "3. Issuing a card..."
echo "Note: Replace USER_TOKEN with the actual token from the registration response"
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d "{
    \"card_type\": \"credit\",
    \"user_token\": \"USER_TOKEN\"
  }"

echo -e "\n\n"

echo "API testing completed!"
