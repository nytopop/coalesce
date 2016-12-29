# TODO for alpha release

## Design changes

### DB migration
MongoDB --> SQLite3

#### Tables
users
	userid     PKEY
	username
	privlevel

posts
	postid     PKEY
	userid     --> users/userid
	title
	body
	bodyHTML
	categoryid --> categories/categoryid
	posted
	updated

categories
	categoryid PKEY
	name

comments
	commentid  PKEY
	postid     --> posts/postid
	parentid   --> comments/commentid || nul
	userid     --> users/userid
	body
	bodyHTML
	posted
	updated

images
	imageid    PKEY
	userid	   --> users/userid
	md5
	thumb_md5
	filename

errors
	errorid    PKEY
	errortext

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

## fixes
random cookie secret

add salt for password tokens

errors are handled incorrectly

comment chains

## features
search function, mongodb full text indexes

spam detection / prevention

asynchronous cortical.io

user profiles

comment deletion

rss feeds

error log in mongo

## deployment lifecycle
build docker host / swarm with ansible

push docker-compose.yml to swarm, with ansible
