# Firebase Database Schema version 1
``` json
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
        "hasAccount": false,  // Will become true when they create a password
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
        "userId": null,  // Initially null, will be updated when profile is created
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
 "contacts": {
      "$submissionId": {
        "timestamp": 1678654321,
        "name": "John Doe",
        "email": "john@example.com",
        "phone": "+1234567890",
        "message": "I have a question about...",
        "contactPreference": "email"
      }
    },
    "newsletters": {
      "$submissionId": {
        "timestamp": 1678654321,
        "email": "john@example.com",
        "phone": "+1234567890", 
        "firstName": "John",
        "termsAccepted": true,
        "gdprAccepted": true,
        "frequency": "weekly"
      }
    },
    "eventNotifications": {
      "$submissionId": {
        "timestamp": 1678654321,
        "email": "john@example.com",
        "phone": "+1234567890",
        "firstName": "John",
        "termsAccepted": true,
        "gdprAccepted": true,
        "frequency": "asNeeded"
      }
    }
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

Key features of this approach:

Progressive user records: Users exist in the database before they have accounts, with hasAccount: false. When they create a password, you'll just update this field.
Preparation for paid subscriptions: The subscriptions.paid node is ready for when you implement paid features.
Unified user identities: Even without authentication, you create user records from form submissions, using email/phone as identifiers.
Submission tracking: All form submissions are still tracked separately for record-keeping.
User-centric model: Each user has a unified profile, even if they've only filled out one form. This supports future personalization.
Preferences section: Consolidates all user preferences in one place for easy retrieval and update.
Subscriptions section: Clear boolean flags for subscription status.
Messages collection: Stores all message content with type indicators.
Indexes: Help you quickly find users by email or phone.

Implementation workflow:

When someone submits any form:

Check if they exist in indexes by email/phone
If not, create a new user record and update indexes
Store the submission with a reference to the user ID


When you implement authentication later:

Add Firebase Authentication
When a user registers, update their existing record with hasAccount: true
No need to migrate data or change schemas


When you implement paid subscriptions:

Just update the subscriptions.paid node for users who subscribe
Use Firebase Authentication roles/claims for access control



This approach gives you the best of both worlds: immediate form handling capabilities with a path to full user accounts and paid subscriptions without schema changes.

Partner relationships:

Added partnerId to the user profile
Added a dedicated partner field in the relationships node
Added partnerEmail to registration submissions for linking
Created a partnerPairs index to track established partnerships


Friend relationships:

Added a friends object to store friend connections
Each friend entry includes since timestamp and status
Added a separate friendRequests collection to manage pending requests


Implementation flows:
For partner signups:

When two people register as partners, check if both exist
Link them via their user IDs in both profiles
Add an entry to the partnerPairs index

For friend connections:

User A sends a friend request to User B
Create an entry in friendRequests.pending
When User B accepts, update both users' relationships.friends objects
Update the request status to "accepted"


For querying:

Find all subscribers to newsletters: users.where(subscriptions.newsletter, equals, true)
Find users by role: users.where(profile.role, equals, "leader")
Find users who came from Facebook: users.where(profile.signupSource, equals, "fb")

For security rules, ensure you restrict access appropriately:

Only authenticated users can modify their own data
Admin-only access to query across all users
Public write access to form submissions endpoints

This structure gives you flexibility for future profiling while maintaining performance for common queries.



This structure allows you to:

Query all users with partners
Find someone's partner quickly
See all of a user's friends
Track friend request status
Maintain the integrity of relationships (both sides of partnerships/friendships)

