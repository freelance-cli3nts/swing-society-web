#!/bin/bash

echo "Testing Contact Form..."
curl -i -X POST http://localhost:8080/api/contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "message": "This is a test message"
  }'

sleep 2  # Wait to avoid rate limiting

echo -e "\n\nTesting Class Inquiry..."
curl -i -X POST http://localhost:8080/api/class \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Student",
    "email": "student@example.com",
    "classType": "Lindy Hop",
    "level": "Beginner",
    "message": "Interested in classes"
  }'

sleep 2  # Wait to avoid rate limiting

echo -e "\n\nTesting Newsletter..."
curl -i -X POST http://localhost:8080/api/newsletter \
  -H "Content-Type: application/json" \
  -d '{
    "email": "news@example.com",
    "name": "Newsletter Sub"
  }'