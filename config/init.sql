create table posts (
  id serial not null unique,
  title varchar(64),
  dateposted date,
  dateupdated date,
  imageurl varchar(128),
  content text,
  summary text,
  keywords text,
  slug varchar(64),
  primary key(id)
);
create table categories (
  id serial not null unique,
  category varchar(32) not null unique,
  primary key (id)
);
-- bridge table between post and category
create table post_categories (
  post_id int not null references posts,
  category_id int not null references categories,
  primary key (post_id, category_id)
);
-- users: ref user_role
create table users (
  id serial not null unique,
  name varchar(28) not null unique,
  email varchar(50) not null unique,
  password text not null,
  primary key (id)
);

insert into posts(title, dateposted, dateupdated, imageurl, content, summary, keywords, slug)
values
  ('Hello World','2020-01-01',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>Hello World</h3>','Hello World','hello world','hello-world'),
  ('The Second Post','2020-01-07',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>A Post</h3>','This is a post','test','second-post'),
  ('The Post','2020-01-14',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>The Post</h3>','The third post','test','more-tests'),
  ('Lets Have another Post','2020-01-14',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>A simple test, hmmm</h3>','May the Forth test - test','test','another-test');

insert into categories(category) values
  ('general'),
  ('off topic'),
  ('news'),
  ('tutorial'),
  ('announcement');

insert into post_categories(post_id, category_id) values
  (1, (select id from categories where id=1) ),
  (1, (select id from categories where id=2) ),
  (2, (select id from categories where id=3) ),
  (3, (select id from categories where id=5) ),
  (4, (select id from categories where id=4) );

