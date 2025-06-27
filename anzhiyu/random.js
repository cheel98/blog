var posts=["posts/a5f1835d/","posts/0/","posts/cf92dd80/"];function toRandomPost(){
    pjax.loadUrl('/'+posts[Math.floor(Math.random() * posts.length)]);
  };