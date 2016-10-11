# coalesce

content for you

[features!]
'feed' based posting system
  your feeds
  topic based
  restricted, public, 
hierarchical comments on each post

TODO
streams/feeds
user pages
full crud of stuff

milestone: commenting system
  threaded comments
  each Post has a Comments document

cortical.io, should generate tags from text, not markup
cortical.io, should use async to generate tags, not as the post happens

error page[x]

user management
  create[x]
  edit[]
  delete[x]

post management
  delete[x]
  edit[x]
  tags[x]
  create[x]
  markdown[x]

comment management
  create[x]
  edit[]
  delete[]


0 - guest
1 - commentor
2 - editor
3 - admin

TODO
  do not allow delete admin
  allow users to post, if verified[x]
  ensure auth works on per user basis, CRUD for own posts only
  make user management from admin, CRUD / verify
  make 'verified' users, which are allowed to post, flag in user db
  build edit/publish/conceal/delete button on /user/me page
  build search function
  change ".List" to ".Comments" and ".Posts"
