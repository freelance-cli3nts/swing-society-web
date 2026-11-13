# Firebase Database Implementation

This document explains how the Firebase Realtime Database integration for the Swing Society Website has been implemented, based on the provided database schema.

## Schema Overview

The database follows a user-centric model with the following main collections:

1. **Users** - Central storage for all user information
   - Profile - Basic user information
   - Auth - Authentication status
   - Subscriptions - Newsletter, event notifications, and paid subscriptions
   - Preferences - User preferences
   - Relationships - Friend and partner connections
   - Messages - User messages from forms

2. **Submissions** - Record of all form submissions
   - Registrations - Class registration forms
   - Contacts - Contact form submissions
   - Newsletters - Newsletter sign-ups
   - EventNotifications - Event notification sign-ups

3. **Indexes** - Quick lookup tables
   - emailToUser - Maps email to user ID
   - phoneToUser - Maps phone to user ID
   - partnerPairs - Tracks partner relationships

## Implementation Details

### User Creation Flow

1. When a form is submitted (registration, contact, or newsletter):
   - Check if the user already exists by email
   - If not, create a new user record
   - Store the form submission in the appropriate submissions collection
   - Link the submission to the user via userId field
   - Update email and phone indexes

2. For existing users:
   - Add new messages to the user's messages collection
   - Update the user's profile if necessary
   - Update subscription status as needed

### Security and Error Handling

- All Firebase operations include error logging
- Operations are structured to be resilient to failures
- In-memory storage provides fallback capability
- User data is protected by appropriate indexing

### Form-Specific Implementation

1. **Registration Form**:
   - Creates user profile with dance role preferences
   - Stores partner information for pair signups
   - Records optional messages

2. **Contact Form**:
   - Creates user profile if needed
   - Adds message to user's messages collection
   - Sets contact preference

3. **Newsletter Subscription**:
   - Creates minimal user profile if needed
   - Sets newsletter subscription status
   - Supports unsubscribe functionality

## Data Storage Path Structure

```
/users/{userId}/
  - profile
  - auth
  - subscriptions
  - preferences
  - relationships
  - messages

/submissions/
  - registrations/{submissionId}
  - contacts/{submissionId}
  - newsletters/{submissionId}
  - eventNotifications/{submissionId}

/indexes/
  - emailToUser/{email} -> userId
  - phoneToUser/{phone} -> userId
  - partnerPairs/{userId1_userId2} -> relationship details
```

## Usage Notes

- Each storage implementation (RegistrationStorage, ContactStorage, NewsletterStorage) has been updated to follow this schema
- The Firebase client provides methods for:
  - Saving form submissions
  - Creating and updating users
  - Managing indexes
  - Querying user information

- The model is designed to support future features including:
  - User authentication
  - Paid subscriptions
  - Friend connections
  - Partner management