# Youtube clone backend

A simple Go based API backend for a Youtube clone. To be implemented as a server based application first and converted to a peer to peer application later on.

## Features to be Implemented

[ ] Upload videos
[ ] Download videos
[ ] Delete videos
[ ] Update video content
[ ] Authentication via email and password
[ ] Social authentication via google, twitter, facebook
[ ] Create Channel
[ ] Delete Channel
[ ] Get list of Channels
[ ] Associate videos to user's channel
[ ] Search functionality
[ ] Follow and followers functionality
[ ] Chat functionality
[ ] Comment functionality
[ ] Like and dislike functionality
[ ] Save to liked videos
[ ] Save to Favourites
[ ] Create playlist
[ ] Add playlist to profile

## Database Tables

- User

  - username
  - email
  - password
  - date_created
  - last_login

- Profile

  - liked videos
  - playlist
  - Followers
  - subscriptions

- Videos

  - \*user
  - video src
  - video title
  - video description
  - likes and dislikes
  - ratings
  - \*comment

- Comments

  - comment
  - \*user

- Channel

  - \*subscriber(s)
  - \*video(s)
  - channel name
  - channel description
  - followers
  - rating

- Subscriber

  - \*user
  - \*channel

- Rating
  - rating count
  - \*user
  - comment

## UML diagram

Coming soon

## Requests screenshots

Coming soon

## Valid URL requests

Coming soon
