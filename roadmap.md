# TODO for alpha release

## Design changes

### DB migration
MongoDB --> SQLite3

#### Priveleges
0 --> guest, read only
1 --> commentor, can CRUD own comments
2 --> editor, can CRUD own comments, own posts
3 --> admin, can CRUD any comments, any posts, any users

### Rewrite HTML / CSS

Text only for now, no images

/auth/sign-in

/auth/register

/config/edit

/error

/img/all

/img/new

/posts/[n]
	/posts redirects to /posts/0 [post 0 - post 30]

	/posts/1 [31 - 60]

	etc etc

/posts/me

/posts/view/[id]

/posts/new

/posts/edit

/users/all

## features
search function, mongodb full text indexes

spam detection / prevention

user profiles

comment deletion

rss feeds
