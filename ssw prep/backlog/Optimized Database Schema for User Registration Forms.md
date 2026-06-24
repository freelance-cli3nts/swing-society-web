The website uses Firebase real time noSQL database. It handles few form submissions: 
1. Registration form
   1. name
   2. email
   3. phone
   4. role: leader / follower / solo
   5. sign up with partner / alone
   6. location of sign-up: fb/google/friends/other
   7. optional message
   8. Terms & Conditions
   9. GDPR
2. Contact form
   1. name
   2. email
   3. phone
   4. message
   5. preference for contacting: email/phone/message app 
3. Newsletter sign-up
   1. email / phone 
   2. first name / nickname
   3. Terms & Conditions
   4. GDPR
   5. Frequency of newsletter
4. Event notification sign-up
   1. email / phone 
   2. first name / nickname
   3. Terms & Conditions
   4. GDPR
   5. Frequency of notification

Key design decisions:
1. **User-centric model**: Each user has a unified profile, even if they've only filled out one form. This supports future personalization
  1. Users exist in the database before they have accounts, with `hasAccount: false`. When they create a password, you'll just update this field.
  2. Creates a profiles collection that builds user profiles based on submitted information
  3. Supports profiling without requiring user accounts
  4. Even without authentication, you create user records from form submissions, using email/phone as identifiers.

For profiling and personalization:
- The `profiles` node will automatically aggregate all information from various submissions
- You can easily look up a person by email or phone
- The profile contains their latest preferences and role information
- You can see which forms they've submitted
- No authentication or accounts needed

2. **Form submissions storage**: 
  1. All form submissions are stored separately, both for record-keeping and to handle cases where someone submits a form without creating an account.
  2. All form submissions are still tracked separately for record-keeping.
  3. Tracks submission IDs to link back to original submissions
  4. When someone submits any form:
   - Check if they exist in indexes by email/phone
   - If not, create a new user record and update indexes
   - Store the submission with a reference to the user ID
3. **Indexes**: These help you quickly find a user by email or phone without querying the entire users collection.
4. **Preferences section**: Consolidates all user preferences in one place for easy retrieval and update.
5. **Subscriptions section**: Clear boolean flags for subscription status. Preparation for paid subscriptions: The `subscriptions.paid` node is ready for when you implement paid features.
  1. Add Firebase Authentication
  2. When a user registers, update their existing record with `hasAccount: true`
  3. No need to migrate data or change schemas
  4. When you implement paid subscriptions:
   - Just update the `subscriptions.paid` node for users who subscribe
   - Use Firebase Authentication roles/claims for access control
6. **Messages collection**: Stores all message content with type indicators.
7. **Partner relationships**:
   - Added `partnerId` to the user profile
   - Added a dedicated `partner` field in the relationships node
   - Added `partnerEmail` to registration submissions for linking
   - Created a `partnerPairs` index to track established partnerships
8. **Friend relationships**:
   - Added a `friends` object to store friend connections
   - Each friend entry includes since timestamp and status
   - Added a separate `friendRequests` collection to manage pending requests
9.  **Implementation flows**:
   For partner signups:
   - When two people register as partners, check if both exist
   - Link them via their user IDs in both profiles
   - Add an entry to the `partnerPairs` index
   For friend connections:
   - User A sends a friend request to User B
   - Create an entry in `friendRequests.pending`
   - When User B accepts, update both users' `relationships.friends` objects
   - Update the request status to "accepted"



For querying:
- Find all newsletter subscribers: `profiles.byEmail.where(profile.newsletterSubscribed, equals, true)`
- Find users by role: `profiles.byEmail.where(profile.role, equals, "leader")`
- Find submissions from a specific source: `submissions.registrations.where(signupSource, equals, "fb")`
- Query all users with partners:
   - Query the `users` collection where `profile.partnerId` field exists and is not null:
     - Example: `users.where('profile.partnerId', '!=', null)`
- Find someone's partner quickly:
   - Given a user ID, retrieve their `profile.partnerId`.  
   - Use this `partnerId` to lookup the corresponding user in the `users` collection.
- See all of a user's friends:
   - Each user profile contains a `relationships.friends` object.  
   - List all friend user IDs found as keys of this object, or iterate over `Object.keys(users.$userId.relationships.friends)`.
- Track friend request status:
   - Use the `friendRequests` collection to view all pending or accepted friend requests.
   - Filter by status field (`pending`, `accepted`, etc.) and/or by sender/recipient user ID.
- Maintain relationship integrity (both sides of partnerships/friendships):
   - When creating or updating a partnership or friendship, ensure both users’ `profile.partnerId` or `relationships.friends` are updated.
   - Always update both corresponding entries (for example, both users’ `relationships.friends` objects when a friendship is created or ended).
   - For partner relationships, update the `partnerPairs` index to reflect established or deleted partnerships.



For security rules, ensure you restrict access appropriately:
- Only authenticated users can modify their own data
- Admin-only access to query across all users
- Public write access to form submissions endpoints

Schema formats: 
1. **JSON Format**: ```schema.json```
2. **Markdown Documentation**: ```database-schema.md```
3. **Schema Dump**: ```schema.txt```


This approach gives you the benefits of user profiling without requiring authentication, while still keeping all your form submission data organized.




```
{
  "users": {
    "$userId": {
      "profile": {
        "name": "John Doe",
        "email": "john@example.com",
        "phone": "+1234567890",
        "role": "follower",
        "signupMethod": "withPartner",
        "signupSource": "friends",
        "partnerId": "user456",  // ID of the partner they signed up with
        "createdAt": 1678654321,
        "updatedAt": 1678654321,
        "termsAccepted": true,
        "gdprAccepted": true
      },
      "auth": {
        "hasAccount": false,
        "accountCreatedAt": null
      },
      "subscriptions": {
        "newsletter": {
          "subscribed": true,
          "frequency": "weekly"
        },
        "eventNotifications": {
          "subscribed": true,
          "frequency": "asNeeded"
        },
        "paid": {
          "active": false,
          "plan": null,
          "startDate": null,
          "endDate": null,
          "autoRenew": false
        }
      },
      "preferences": {
        "contactPreference": "email"
      },
      "relationships": {
        "partner": "user456",  // Direct reference to partner
        "friends": {
          "user789": {
            "since": 1678654321,
            "status": "confirmed"  // confirmed, pending, requested
          },
          "user101": {
            "since": 1678700000,
            "status": "pending"
          }
        }
      },
      "messages": {
        "$messageId": {
          "content": "I'd like more information about...",
          "timestamp": 1678654321,
          "type": "registration"
        }
      }
    }
  },
  "submissions": {
    "registrations": {
      "$submissionId": {
        "userId": null,
        "partnerId": null,  // For tracking partner signups
        "timestamp": 1678654321,
        "name": "John Doe",
        "email": "john@example.com",
        "phone": "+1234567890",
        "role": "follower",
        "signupMethod": "withPartner",
        "partnerEmail": "partner@example.com",  // To link partners during registration
        "signupSource": "friends",
        "message": "I'd like more information about...",
        "termsAccepted": true,
        "gdprAccepted": true
      }
    },
    "contacts": { /* Same structure as before */ },
    "newsletters": { /* Same structure as before */ },
    "eventNotifications": { /* Same structure as before */ }
  },
  "indexes": {
    "emailToUser": {
      "john@example.com": "user123"
    },
    "phoneToUser": {
      "+1234567890": "user123"
    },
    "partnerPairs": {
      "user123_user456": {
        "timestamp": 1678654321,
        "active": true
      }
    }
  },
  "friendRequests": {
    "pending": {
      "$requestId": {
        "from": "user123",
        "to": "user789",
        "timestamp": 1678654321,
        "message": "Let's connect!",
        "status": "pending"  // pending, accepted, declined
      }
    }
  }
}
```



