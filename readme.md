# coalesce

Lightning fast blogging CMS

[go]
auth.go middleware
form validation (auth, write)
error.go

[html]
sign in page / form
client side form validation
Recent Posts tooltip
index page

[editable elements]
Users
Pages
Posts
Comments

finish auth scheme
milestone: pages
milestone: commenting system
  threaded comments
  each Post has a Comments document,
  

need to expand the json structure in template format

range through .Comments
  range through .Replies
	range through .Replies
	  recurse until empty
