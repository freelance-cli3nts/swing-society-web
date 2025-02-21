# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting API Tests...${NC}\n"

# Test Contact Form
echo -e "${GREEN}Testing Contact Form...${NC}"
curl -i -X POST http://localhost:3001/api/contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "message": "This is a test message"
  }'

sleep 1

# Test Class Inquiry
echo -e "\n\n${GREEN}Testing Class Inquiry...${NC}"
curl -i -X POST http://localhost:3001/api/class \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Student",
    "email": "student@example.com",
    "classType": "Lindy Hop",
    "level": "Beginner",
    "message": "Interested in classes"
  }'

sleep 1

# Test Newsletter
echo -e "\n\n${GREEN}Testing Newsletter...${NC}"
curl -i -X POST http://localhost:3001/api/newsletter \
  -H "Content-Type: application/json" \
  -d '{
    "email": "news@example.com",
    "name": "Newsletter Sub"
  }'

sleep 1

# Test Registration
echo -e "\n\n${GREEN}Testing Registration...${NC}"
curl -i -X POST http://localhost:3001/api/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=Test+Registration&email=register@example.com"

sleep 1

# Test Carousel
echo -e "\n\n${GREEN}Testing Carousel Data...${NC}"
curl -i http://localhost:3001/api/carousel/general

sleep 1

# Test Template (with HTMX header)
echo -e "\n\n${GREEN}Testing Template with HTMX...${NC}"
curl -i http://localhost:3001/templates/contact \
  -H "HX-Request: true"

echo -e "\n\n${GREEN}All tests completed!${NC}"